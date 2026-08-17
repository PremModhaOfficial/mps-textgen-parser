# Go TextGen Reference — Full Dynamic Generation

## 1. Infra Concept (add to Structure aspect)

Root concept, implements INamedConcept, alias: `infra`

**Properties:**
- `modulePath` : string  — e.g. "github.com/user/project"
- `port` : string — e.g. "8080"
- `dbUser` : string — e.g. "postgres"
- `dbPass` : string — e.g. "password"
- `dbName` : string — e.g. "mydb"

---

## 2. Code Concept (add to Structure aspect)

Root concept, implements INamedConcept, alias: `code`

**References:**
- `model` : Models (cardinality 1)
- `infra` : Infra (cardinality 1)

---

## 3. Behavior Methods (add in Behavior aspect)

### On ModelSchema concept:

#### `pascalCase(string input)` — static method, returns string
```java
StringBuilder sb = new StringBuilder();
string[] parts = input.split("_");
for (int i = 0; i < parts.length; i++) {
  string part = parts[i];
  if (part.length() > 0) {
    sb.append(part.substring(0, 1).toUpperCase());
    if (part.length() > 1) {
      sb.append(part.substring(1));
    }
  }
}
return sb.toString();
```

#### `singularName()` — returns string
```java
string name = this.name;
if (name.endsWith("s") && name.length() > 1) {
  return name.substring(0, name.length() - 1);
}
return name;
```

#### `structName()` — returns string
```java
return ModelSchema.pascalCase(this.singularName());
```

#### `repoName()` — returns string
```java
return this.structName() + "Repo";
```

#### `createStructName()` — returns string
```java
return this.structName() + "Create";
```

#### `hasReferences()` — returns boolean
```java
for (node<Field_PlaceHolder> f : this.Fields) {
  if (f.isInstanceOf(FieldRefrence)) {
    return true;
  }
}
return false;
```

### On Field concept:

#### `pascalName()` — returns string
```java
return ModelSchema.pascalCase(this.name);
```

#### `goType()` — returns string
```java
if (this.dataType == DataType.int64) { return "int64"; }
if (this.dataType == DataType.string) { return "string"; }
if (this.dataType == DataType.float) { return "float64"; }
if (this.dataType == DataType.boolean) { return "bool"; }
if (this.dataType == DataType.timestamp) { return "time.Time"; }
return "interface{}";
```

### On FieldRefrence concept:

#### `pascalName()` — returns string
```java
return ModelSchema.pascalCase(this.name);
```

---

## 4. TextGen for Code concept

```
text gen component for concept Code {
  file name : (node)->string {
    node.model.name;
  }
  (node)->void {
    node<Code> n = node;
    node<Models> model = n.model;
    node<Infra> infra = n.infra;
```

### SECTION A — Package + Imports

```
    append {package main\n\n} ;
    append {import (\n} ;
    append {\t"database/sql"\n} ;
    append {\t_ "embed"\n} ;
    append {\t"encoding/json"\n} ;
    append {\t"fmt"\n} ;
    append {\t"log"\n} ;
    append {\t"net/http"\n} ;
    append {\t"os"\n} ;
    append {\t"strconv"\n} ;
    append {\t"time"\n\n} ;
    append {\t_ "github.com/lib/pq"\n} ;
    append {\thttpSwagger "github.com/swaggo/http-swagger"\n} ;
    append {\t_ "} ${infra.modulePath} {/docs"\n} ;
    append {)\n\n} ;
```

### SECTION B — Embed + Swagger

```
    append {//go:embed user_management_init.sql\n} ;
    append {var migrationSQL string\n\n} ;
    append {// @title         } ${model.name} { API\n} ;
    append {// @version       1.0\n} ;
    append {// @description   CRUD service for } ${model.name} {\n} ;
    append {// @host          localhost:} ${infra.port} {\n} ;
    append {// @BasePath      /\n} ;
    append {// @schemes       http\n} ;
    append {// @produce       json\n} ;
    append {// @consumes      json\n\n} ;
```

### SECTION C — Model Structs (regular schemas)

```
    append {// ============================================================\n} ;
    append {// Models\n} ;
    append {// ============================================================\n\n} ;

    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {

        // ---- Main struct ----
        append {type } ${schema.structName()} { struct {\n} ;
        append {\tID int64 `json:"id" db:"id" example:"1"`\n} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {\t} ${f.pascalName()} { } ${f.goType()} { `json:"} ${f.name} {" db:"} ${f.name} {"`\n} ;
          }
        }
        append {}\n\n} ;

        // ---- MarshalJSON (if any Sensitive field) ----
        boolean hasSensitive = false;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.Sensitive) {
              hasSensitive = true;
            }
          }
        }
        if (hasSensitive) {
          append {func (u } ${schema.structName()} {) MarshalJSON() ([]byte, error) {\n} ;
          append {\ttype Alias } ${schema.structName()} {\n} ;
          append {\treturn json.Marshal(&struct {\n} ;
          append {\t\tAlias\n} ;
          foreach field in schema.Fields {
            if (field.isInstanceOf(Field)) {
              node<Field> f = (Field) field;
              if (f.Sensitive) {
                append {\t\t} ${f.pascalName()} { string `json:"} ${f.name} {"`\n} ;
              }
            }
          }
          append {\t}{\n} ;
          append {\t\tAlias: (Alias)(u),\n} ;
          foreach field in schema.Fields {
            if (field.isInstanceOf(Field)) {
              node<Field> f = (Field) field;
              if (f.Sensitive) {
                append {\t\t} ${f.pascalName()} {: "[REDACTED]",\n} ;
              }
            }
          }
          append {\t})\n}\n\n} ;
        }

        // ---- Create struct (skip timestamp fields) ----
        append {type } ${schema.createStructName()} { struct {\n} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType != DataType.timestamp) {
              append {\t} ${f.pascalName()} { } ${f.goType()} { `json:"} ${f.name} {"`\n} ;
            }
          }
        }
        append {}\n\n} ;
      }
    }
```

### SECTION D — Model Structs (join schemas)

```
    foreach schema in model.schemas {
      if (schema.hasReferences()) {

        // ---- Join struct ----
        append {type } ${schema.structName()} { struct {\n} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            append {\t} ${fr.pascalName()} { int64 `json:"} ${fr.name} {" db:"} ${fr.name} {"`\n} ;
          }
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {\t} ${f.pascalName()} { } ${f.goType()} { `json:"} ${f.name} {" db:"} ${f.name} {"`\n} ;
          }
        }
        append {}\n\n} ;

        // ---- Assign body struct (first FK only) ----
        // Get the second reference for the body
        int refCount = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (refCount == 1) {
              append {type Assign} ${fr.target_schema.structName()} {Body struct {\n} ;
              append {\t} ${fr.pascalName()} { int64 `json:"} ${fr.name} {" binding:"required"`\n} ;
              append {}\n\n} ;
            }
            refCount = refCount + 1;
          }
        }
      }
    }
```

### SECTION E — Repositories (regular schemas)

```
    append {// ============================================================\n} ;
    append {// Repositories\n} ;
    append {// ============================================================\n\n} ;

    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {
        string sn = schema.structName();
        string rn = schema.repoName();
        string tn = schema.name;

        append {type } ${rn} { struct{ db *sql.DB }\n\n} ;

        // ---- Create ----
        append {func (r *} ${rn} {) Create(u *} ${sn} {) error {\n} ;
        append {\treturn r.db.QueryRow(\n} ;
        append {\t\t`INSERT INTO } ${tn} { (} ;
        int idx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType != DataType.timestamp) {
              if (idx > 0) { append {, } ; }
              append ${f.name} ;
              idx = idx + 1;
            }
          }
        }
        append {)\n\t\t VALUES (} ;
        for (int i = 1; i <= idx; i++) {
          if (i > 1) { append {, } ; }
          append {$} ${i} ;
        }
        append {)\n\t\t RETURNING id} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType == DataType.timestamp) {
              append {, } ${f.name} ;
            }
          }
        }
        append {`,\n} ;

        // args
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType != DataType.timestamp) {
              append {\t\tu.} ${f.pascalName()} {,\n} ;
            }
          }
        }
        append {\t).Scan(&u.ID} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType == DataType.timestamp) {
              append {, &u.} ${f.pascalName()} ;
            }
          }
        }
        append {)\n}\n\n} ;

        // ---- GetByID ----
        append {func (r *} ${rn} {) GetByID(id int64) (*} ${sn} {, error) {\n} ;
        append {\tu := &} ${sn} {{}}\n} ;
        append {\terr := r.db.QueryRow(\n} ;
        append {\t\t`SELECT id} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, } ${f.name} ;
          }
        }
        append {\n\t\t FROM } ${tn} { WHERE id = $1`, id,\n} ;
        append {\t).Scan(&u.ID} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, &u.} ${f.pascalName()} ;
          }
        }
        append {)\n} ;
        append {\treturn u, err\n}\n\n} ;

        // ---- List ----
        append {func (r *} ${rn} {) List() ([]} ${sn} {, error) {\n} ;
        append {\trows, err := r.db.Query(\n} ;
        append {\t\t`SELECT id} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, } ${f.name} ;
          }
        }
        append {\n\t\t FROM } ${tn} { ORDER BY id`)\n} ;
        append {\tif err != nil {\n\t\treturn nil, err\n\t}\n} ;
        append {\tdefer rows.Close()\n} ;
        append {\tvar items []} ${sn} {\n} ;
        append {\tfor rows.Next() {\n} ;
        append {\t\tvar u } ${sn} {\n} ;
        append {\t\tif err := rows.Scan(&u.ID} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, &u.} ${f.pascalName()} ;
          }
        }
        append {); err != nil {\n\t\t\treturn nil, err\n\t\t}\n} ;
        append {\t\titems = append(items, u)\n\t}\n} ;
        append {\treturn items, rows.Err()\n}\n\n} ;

        // ---- Update ----
        append {func (r *} ${rn} {) Update(u *} ${sn} {) error {\n} ;
        append {\treturn r.db.QueryRow(\n} ;
        append {\t\t`UPDATE } ${tn} { SET } ;
        int uidx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType != DataType.timestamp) {
              if (uidx > 0) { append {, } ; }
              uidx = uidx + 1;
              append ${f.name} { = $} ${uidx} ;
            }
          }
        }
        int nextParam = uidx + 1;
        append {, updated_at = NOW()\n\t\t WHERE id = $} ${nextParam} ;
        append {\n\t\t RETURNING updated_at`,\n} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType != DataType.timestamp) {
              append {\t\tu.} ${f.pascalName()} {,\n} ;
            }
          }
        }
        append {\t\tu.ID,\n} ;
        append {\t).Scan(&u.UpdatedAt)\n}\n\n} ;

        // ---- Delete ----
        append {func (r *} ${rn} {) Delete(id int64) error {\n} ;
        append {\t_, err := r.db.Exec(`DELETE FROM } ${tn} { WHERE id = $1`, id)\n} ;
        append {\treturn err\n}\n\n} ;
      }
    }
```

### SECTION F — Repositories (join schemas)

```
    foreach schema in model.schemas {
      if (schema.hasReferences()) {
        string sn = schema.structName();
        string rn = schema.repoName();
        string tn = schema.name;

        append {type } ${rn} { struct{ db *sql.DB }\n\n} ;

        // ---- Assign ----
        append {func (r *} ${rn} {) Assign(} ;
        int argIdx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (argIdx > 0) { append {, } ; }
            append ${fr.name} { int64} ;
            argIdx = argIdx + 1;
          }
        }
        append {) (*} ${sn} {, error) {\n} ;
        append {\tur := &} ${sn} {{}}\n} ;
        append {\terr := r.db.QueryRow(\n} ;
        append {\t\t`INSERT INTO } ${tn} { (} ;
        int fkIdx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (fkIdx > 0) { append {, } ; }
            append ${fr.name} ;
            fkIdx = fkIdx + 1;
          }
        }
        append {) VALUES (} ;
        for (int i = 1; i <= fkIdx; i++) {
          if (i > 1) { append {, } ; }
          append {$} ${i} ;
        }
        append {)\n\t\t ON CONFLICT (} ;
        int ckIdx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (ckIdx > 0) { append {, } ; }
            append ${fr.name} ;
            ckIdx = ckIdx + 1;
          }
        }
        append {) DO NOTHING\n\t\t RETURNING } ;
        int retIdx = 0;
        foreach field in schema.Fields {
          if (retIdx > 0) { append {, } ; }
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            append ${fr.name} ;
          }
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append ${f.name} ;
          }
          retIdx = retIdx + 1;
        }
        append {`,\n} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            append {\t\t} ${fr.name} {,\n} ;
          }
        }
        append {\t).Scan(} ;
        int scanIdx = 0;
        foreach field in schema.Fields {
          if (scanIdx > 0) { append {, } ; }
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            append {&ur.} ${fr.pascalName()} ;
          }
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {&ur.} ${f.pascalName()} ;
          }
          scanIdx = scanIdx + 1;
        }
        append {)\n} ;
        append {\treturn ur, err\n}\n\n} ;

        // ---- Remove ----
        append {func (r *} ${rn} {) Remove(} ;
        int rmArgIdx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (rmArgIdx > 0) { append {, } ; }
            append ${fr.name} { int64} ;
            rmArgIdx = rmArgIdx + 1;
          }
        }
        append {) error {\n} ;
        append {\t_, err := r.db.Exec(`DELETE FROM } ${tn} { WHERE } ;
        int whereIdx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (whereIdx > 0) { append { AND } ; }
            whereIdx = whereIdx + 1;
            append ${fr.name} { = $} ${whereIdx} ;
          }
        }
        append {`} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            append {, } ${fr.name} ;
          }
        }
        append {)\n\treturn err\n}\n\n} ;

        // ---- Cross-query: for each ref A, get target_schema(B) by A ----
        // e.g. GetRolesByUser, GetUsersByRole
        foreach refA in schema.Fields {
          if (refA.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> frA = (FieldRefrence) refA;
            foreach refB in schema.Fields {
              if (refB.isInstanceOf(FieldRefrence) && refB != refA) {
                node<FieldRefrence> frB = (FieldRefrence) refB;
                // Generate: Get <B.target plural> By <A.target singular>
                // e.g. GetRolesByUser(userId int64) ([]Role, error)
                string targetStruct = frB.target_schema.structName();
                string targetTable = frB.target_schema.name;

                append {func (r *} ${rn} {) Get} ${frB.target_schema.structName()} {sBy} ${frA.target_schema.structName()} {(} ${frA.name} { int64) ([]} ${targetStruct} {, error) {\n} ;
                append {\trows, err := r.db.Query(\n} ;
                append {\t\t`SELECT t.id} ;
                foreach tf in frB.target_schema.Fields {
                  if (tf.isInstanceOf(Field)) {
                    node<Field> f = (Field) tf;
                    append {, t.} ${f.name} ;
                  }
                }
                append {\n\t\t FROM } ${targetTable} { t\n} ;
                append {\t\t INNER JOIN } ${tn} { j ON j.} ${frB.name} { = t.id\n} ;
                append {\t\t WHERE j.} ${frA.name} { = $1\n} ;
                append {\t\t ORDER BY t.id`, } ${frA.name} {,\n} ;
                append {\t)\n} ;
                append {\tif err != nil {\n\t\treturn nil, err\n\t}\n} ;
                append {\tdefer rows.Close()\n} ;
                append {\tvar items []} ${targetStruct} {\n} ;
                append {\tfor rows.Next() {\n} ;
                append {\t\tvar item } ${targetStruct} {\n} ;
                append {\t\tif err := rows.Scan(&item.ID} ;
                foreach tf in frB.target_schema.Fields {
                  if (tf.isInstanceOf(Field)) {
                    node<Field> f = (Field) tf;
                    append {, &item.} ${f.pascalName()} ;
                  }
                }
                append {); err != nil {\n\t\t\treturn nil, err\n\t\t}\n} ;
                append {\t\titems = append(items, item)\n\t}\n} ;
                append {\treturn items, rows.Err()\n}\n\n} ;
              }
            }
          }
        }
      }
    }
```

### SECTION G — HTTP Handlers (regular schemas)

```
    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {
        string sn = schema.structName();
        string rn = schema.repoName();
        string tn = schema.name;

        append {// ============================================================\n} ;
        append {// HTTP Handlers — } ${sn} {\n} ;
        append {// ============================================================\n\n} ;

        // ---- handleCreate ----
        append {func handleCreate} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n} ;
        append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
        append {\t\tvar u } ${sn} {\n} ;
        append {\t\tif err := json.NewDecoder(r.Body).Decode(&u); err != nil {\n} ;
        append {\t\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tif err := repo.Create(&u); err != nil {\n} ;
        append {\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tw.Header().Set("Content-Type", "application/json")\n} ;
        append {\t\tw.WriteHeader(http.StatusCreated)\n} ;
        append {\t\tjson.NewEncoder(w).Encode(u)\n} ;
        append {\t}\n}\n\n} ;

        // ---- handleGet ----
        append {func handleGet} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n} ;
        append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
        append {\t\tidStr := r.PathValue("id")\n} ;
        append {\t\tid, err := strconv.ParseInt(idStr, 10, 64)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tu, err := repo.GetByID(id)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusNotFound)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tw.Header().Set("Content-Type", "application/json")\n} ;
        append {\t\tjson.NewEncoder(w).Encode(u)\n} ;
        append {\t}\n}\n\n} ;

        // ---- handleList ----
        append {func handleList} ${sn} {s(repo *} ${rn} {) http.HandlerFunc {\n} ;
        append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
        append {\t\titems, err := repo.List()\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tw.Header().Set("Content-Type", "application/json")\n} ;
        append {\t\tjson.NewEncoder(w).Encode(items)\n} ;
        append {\t}\n}\n\n} ;

        // ---- handleUpdate ----
        append {func handleUpdate} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n} ;
        append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
        append {\t\tidStr := r.PathValue("id")\n} ;
        append {\t\tid, err := strconv.ParseInt(idStr, 10, 64)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tvar u } ${sn} {\n} ;
        append {\t\tif err := json.NewDecoder(r.Body).Decode(&u); err != nil {\n} ;
        append {\t\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tu.ID = id\n} ;
        append {\t\tif err := repo.Update(&u); err != nil {\n} ;
        append {\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tw.Header().Set("Content-Type", "application/json")\n} ;
        append {\t\tjson.NewEncoder(w).Encode(u)\n} ;
        append {\t}\n}\n\n} ;

        // ---- handleDelete ----
        append {func handleDelete} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n} ;
        append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
        append {\t\tidStr := r.PathValue("id")\n} ;
        append {\t\tid, err := strconv.ParseInt(idStr, 10, 64)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tif err := repo.Delete(id); err != nil {\n} ;
        append {\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tw.WriteHeader(http.StatusNoContent)\n} ;
        append {\t}\n}\n\n} ;
      }
    }
```

### SECTION H — HTTP Handlers (join schemas)

```
    foreach schema in model.schemas {
      if (schema.hasReferences()) {
        string sn = schema.structName();
        string rn = schema.repoName();
        string tn = schema.name;

        append {// ============================================================\n} ;
        append {// HTTP Handlers — } ${sn} { (assignments)\n} ;
        append {// ============================================================\n\n} ;

        // Get the two refs
        node<FieldRefrence> firstRef = null;
        node<FieldRefrence> secondRef = null;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (firstRef == null) { firstRef = fr; }
            else if (secondRef == null) { secondRef = fr; }
          }
        }

        // ---- handleAssign ----
        // POST /firstRef.target_schema.name/{id}/secondRef.target_schema.name
        append {func handleAssign} ${secondRef.target_schema.structName()} {(urRepo *} ${rn} {) http.HandlerFunc {\n} ;
        append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
        append {\t\t} ${firstRef.name} {, err := strconv.ParseInt(r.PathValue("id"), 10, 64)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tvar body Assign} ${secondRef.target_schema.structName()} {Body\n} ;
        append {\t\tif err := json.NewDecoder(r.Body).Decode(&body); err != nil {\n} ;
        append {\t\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tur, err := urRepo.Assign(} ${firstRef.name} {, body.} ${secondRef.pascalName()} {)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tw.Header().Set("Content-Type", "application/json")\n} ;
        append {\t\tw.WriteHeader(http.StatusCreated)\n} ;
        append {\t\tjson.NewEncoder(w).Encode(ur)\n} ;
        append {\t}\n}\n\n} ;

        // ---- handleRemove ----
        append {func handleRemove} ${secondRef.target_schema.structName()} {(urRepo *} ${rn} {) http.HandlerFunc {\n} ;
        append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
        append {\t\t} ${firstRef.name} {, err := strconv.ParseInt(r.PathValue("id"), 10, 64)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\t} ${secondRef.name} {, err := strconv.ParseInt(r.PathValue("} ${secondRef.name} {"), 10, 64)\n} ;
        append {\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tif err := urRepo.Remove(} ${firstRef.name} {, } ${secondRef.name} {); err != nil {\n} ;
        append {\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n} ;
        append {\t\tw.WriteHeader(http.StatusNoContent)\n} ;
        append {\t}\n}\n\n} ;

        // ---- handleGet<B>s for each ref pair ----
        foreach refA in schema.Fields {
          if (refA.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> frA = (FieldRefrence) refA;
            foreach refB in schema.Fields {
              if (refB.isInstanceOf(FieldRefrence) && refB != refA) {
                node<FieldRefrence> frB = (FieldRefrence) refB;
                append {func handleGet} ${frA.target_schema.structName()} ${frB.target_schema.structName()} {s(urRepo *} ${rn} {) http.HandlerFunc {\n} ;
                append {\treturn func(w http.ResponseWriter, r *http.Request) {\n} ;
                append {\t\tid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)\n} ;
                append {\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n} ;
                append {\t\titems, err := urRepo.Get} ${frB.target_schema.structName()} {sBy} ${frA.target_schema.structName()} {(id)\n} ;
                append {\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n} ;
                append {\t\tw.Header().Set("Content-Type", "application/json")\n} ;
                append {\t\tjson.NewEncoder(w).Encode(items)\n} ;
                append {\t}\n}\n\n} ;
              }
            }
          }
        }
      }
    }
```

### SECTION I — Main Function

```
    append {// ============================================================\n} ;
    append {// Main\n} ;
    append {// ============================================================\n\n} ;

    append {func main() {\n} ;
    append {\tdbURL := os.Getenv("DATABASE_URL")\n} ;
    append {\tif dbURL == "" {\n} ;
    append {\t\tdbURL = "postgres://} ${infra.dbUser} {:} ${infra.dbPass} {@localhost:5432/} ${infra.dbName} {?sslmode=disable"\n} ;
    append {\t}\n\n} ;

    append {\tdb, err := sql.Open("postgres", dbURL)\n} ;
    append {\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n} ;
    append {\tdefer db.Close()\n\n} ;

    append {\tfor i := 0; i < 5; i++ {\n} ;
    append {\t\tif err = db.Ping(); err == nil {\n\t\t\tbreak\n\t\t}\n} ;
    append {\t\tlog.Printf("DB not ready, retrying... (%d/5)", i+1)\n} ;
    append {\t\ttime.Sleep(2 * time.Second)\n} ;
    append {\t}\n} ;
    append {\tif err != nil {\n\t\tlog.Fatal("DB connection failed:", err)\n\t}\n\n} ;

    append {\tif _, err := db.Exec(migrationSQL); err != nil {\n\t\tlog.Fatal(err)\n\t}\n} ;
    append {\tlog.Println("Migration complete")\n\n} ;

    // Repo instantiation
    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {
        append {\t} ${schema.singularName()} {Repo := &} ${schema.repoName()} {{db: db}\n} ;
      }
    }
    foreach schema in model.schemas {
      if (schema.hasReferences()) {
        append {\t} ${schema.singularName()} {Repo := &} ${schema.repoName()} {{db: db}\n} ;
      }
    }

    append {\n\tmux := http.NewServeMux()\n\n} ;
    append {\t// Swagger UI\n} ;
    append {\tmux.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)\n\n} ;

    // Regular schema routes
    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {
        string varRepo = schema.singularName() + "Repo";
        append {\t// } ${schema.structName()} {s\n} ;
        append {\tmux.HandleFunc("POST /} ${schema.name} {", handleCreate} ${schema.structName()} {(} ${varRepo} {))\n} ;
        append {\tmux.HandleFunc("GET /} ${schema.name} {", handleList} ${schema.structName()} {s(} ${varRepo} {))\n} ;
        append {\tmux.HandleFunc("GET /} ${schema.name} {/{id}", handleGet} ${schema.structName()} {(} ${varRepo} {))\n} ;
        append {\tmux.HandleFunc("PUT /} ${schema.name} {/{id}", handleUpdate} ${schema.structName()} {(} ${varRepo} {))\n} ;
        append {\tmux.HandleFunc("DELETE /} ${schema.name} {/{id}", handleDelete} ${schema.structName()} {(} ${varRepo} {))\n\n} ;
      }
    }

    // Join schema routes
    foreach schema in model.schemas {
      if (schema.hasReferences()) {
        string varRepo = schema.singularName() + "Repo";
        node<FieldRefrence> fRef = null;
        node<FieldRefrence> sRef = null;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (fRef == null) { fRef = fr; }
            else if (sRef == null) { sRef = fr; }
          }
        }

        append {\t// } ${schema.structName()} { assignments\n} ;
        append {\tmux.HandleFunc("POST /} ${fRef.target_schema.name} {/{id}/} ${sRef.target_schema.name} {", handleAssign} ${sRef.target_schema.structName()} {(} ${varRepo} {))\n} ;
        append {\tmux.HandleFunc("DELETE /} ${fRef.target_schema.name} {/{id}/} ${sRef.target_schema.name} {/{} ${sRef.name} {}", handleRemove} ${sRef.target_schema.structName()} {(} ${varRepo} {))\n} ;
        append {\tmux.HandleFunc("GET /} ${fRef.target_schema.name} {/{id}/} ${sRef.target_schema.name} {", handleGet} ${fRef.target_schema.structName()} ${sRef.target_schema.structName()} {s(} ${varRepo} {))\n} ;
        append {\tmux.HandleFunc("GET /} ${sRef.target_schema.name} {/{id}/} ${fRef.target_schema.name} {", handleGet} ${sRef.target_schema.structName()} ${fRef.target_schema.structName()} {s(} ${varRepo} {))\n\n} ;
      }
    }

    append {\tfmt.Println("Serving on :} ${infra.port} {")\n} ;
    append {\tfmt.Println("Swagger UI: http://localhost:} ${infra.port} {/swagger/index.html")\n} ;
    append {\tlog.Fatal(http.ListenAndServe(":} ${infra.port} {", mux))\n} ;
    append {}\n} ;
```

### CLOSING

```
  }
  file path : (node)->string {
    "go/";
  }
  extension : (node)->string {
    "go";
  }
}
```
