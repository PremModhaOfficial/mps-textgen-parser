# Go TextGen — Minimal Appends Version

## New Concepts Needed

### Infra (root, INamedConcept, alias: `infra`)
Props: `modulePath` (string), `port` (string), `dbUser` (string), `dbPass` (string), `dbName` (string)

### Code (root, INamedConcept, alias: `code`)
Refs: `model` → Models (1), `infra` → Infra (1)

---

## Behavior Methods

### On ModelSchema:

**static pascalCase(string input) → string**
```java
StringBuilder sb = new StringBuilder();
string[] parts = input.split("_");
for (int i = 0; i < parts.length; i++) {
  string part = parts[i];
  if (part.length() > 0) {
    sb.append(part.substring(0, 1).toUpperCase());
    if (part.length() > 1) { sb.append(part.substring(1)); }
  }
}
return sb.toString();
```

**singularName() → string**
```java
string name = this.name;
if (name.endsWith("s") && name.length() > 1) { return name.substring(0, name.length() - 1); }
return name;
```

**structName() → string** — `return ModelSchema.pascalCase(this.singularName());`
**repoName() → string** — `return this.structName() + "Repo";`
**createStructName() → string** — `return this.structName() + "Create";`

**hasReferences() → boolean**
```java
for (node<Field_PlaceHolder> f : this.Fields) {
  if (f.isInstanceOf(FieldRefrence)) { return true; }
}
return false;
```

### On Field:
**pascalName() → string** — `return ModelSchema.pascalCase(this.name);`
**goType() → string**
```java
if (this.dataType == DataType.int64) { return "int64"; }
if (this.dataType == DataType.string) { return "string"; }
if (this.dataType == DataType.float) { return "float64"; }
if (this.dataType == DataType.boolean) { return "bool"; }
if (this.dataType == DataType.timestamp) { return "time.Time"; }
return "interface{}";
```

### On FieldRefrence:
**pascalName() → string** — `return ModelSchema.pascalCase(this.name);`

---

## TextGen for Code

```
text gen component for concept Code {
  file name : (node)->string { node.model.name; }
  (node)->void {
    node<Code> n = node;
    node<Models> model = n.model;
    node<Infra> infra = n.infra;

    append {package main\n\nimport (\n\t"database/sql"\n\t_ "embed"\n\t"encoding/json"\n\t"fmt"\n\t"log"\n\t"net/http"\n\t"os"\n\t"strconv"\n\t"time"\n\n\t_ "github.com/lib/pq"\n\thttpSwagger "github.com/swaggo/http-swagger"\n\t_ "} ${infra.modulePath} {/docs"\n)\n\n//go:embed user_management_init.sql\nvar migrationSQL string\n\n// @title         } ${model.name} { API\n// @version       1.0\n// @description   CRUD service for } ${model.name} {\n// @host          localhost:} ${infra.port} {\n// @BasePath      /\n// @schemes       http\n// @produce       json\n// @consumes      json\n\n// ============================================================\n// Models\n// ============================================================\n\n} ;

    // ---- Regular schema structs ----
    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {

        append {type } ${schema.structName()} { struct {\n\tID int64 `json:"id" db:"id" example:"1"`\n} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {\t} ${f.pascalName()} { } ${f.goType()} { `json:"} ${f.name} {" db:"} ${f.name} {"`\n} ;
          }
        }
        append {}\n\n} ;

        // MarshalJSON if Sensitive
        boolean hasSensitive = false;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.Sensitive) { hasSensitive = true; }
          }
        }

        // found the sensitive parst
        if (hasSensitive) {
          append {func (u } ${schema.structName()} {) MarshalJSON() ([]byte, error) {\n\ttype Alias } ${schema.structName()} {\n\treturn json.Marshal(&struct {\n\t\tAlias\n} ;
          foreach field in schema.Fields {
            if (field.isInstanceOf(Field)) {
              node<Field> f = (Field) field;
              if (f.Sensitive) {
                append {\t\t} ${f.pascalName()} { string `json:"} ${f.name} {"`\n} ;
              }
            }
          }
          append {\t}{\n\t\tAlias: (Alias)(u),\n} ;
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


        // Create struct
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

    // ---- Join schema structs ----
    foreach schema in model.schemas {
      if (schema.hasReferences()) {
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

        int refCount = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (refCount == 1) {
              append {type Assign} ${fr.target_schema.structName()} {Body struct {\n\t} ${fr.pascalName()} { int64 `json:"} ${fr.name} {" binding:"required"`\n}\n\n} ;
            }
            refCount = refCount + 1;
          }
        }
      }
    }

    // ============================================================
    // Repositories — regular
    // ============================================================

    append {// ============================================================\n// Repositories\n// ============================================================\n\n} ;

    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {
        string sn = schema.structName();
        string rn = schema.repoName();
        string tn = schema.name;

        append {type } ${rn} { struct{ db *sql.DB }\n\n} ;

        // Create
        append {func (r *} ${rn} {) Create(u *} ${sn} {) error {\n\treturn r.db.QueryRow(\n\t\t`INSERT INTO } ${tn} { (} ;
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
            if (f.dataType == DataType.timestamp) { append {, } ${f.name} ; }
          }
        }
        append {`,\n} ;

        // non timestapm 
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType != DataType.timestamp) { append {\t\tu.} ${f.pascalName()} {,\n} ; }
          }
        }

        // scanning
        append {\t).Scan(&u.ID} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType == DataType.timestamp) { append {, &u.} ${f.pascalName()} ; }
          }
        }
        append {)\n}\n\n} ;

        // GetByID
        append {func (r *} ${rn} {) GetByID(id int64) (*} ${sn} {, error) {\n\tu := &} ${sn} {{}\n\terr := r.db.QueryRow(\n\t\t`SELECT id} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, } ${f.name} ;
          }
        }
        append {\n\t\t FROM } ${tn} { WHERE id = $1`, id,\n\t).Scan(&u.ID} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, &u.} ${f.pascalName()} ;
          }
        }
        append {)\n\treturn u, err\n}\n\n} ;

        // List
        append {func (r *} ${rn} {) List() ([]} ${sn} {, error) {\n\trows, err := r.db.Query(\n\t\t`SELECT id} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, } ${f.name} ;
          }
        }
        append {\n\t\t FROM } ${tn} { ORDER BY id`)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tdefer rows.Close()\n\tvar items []} ${sn} {\n\tfor rows.Next() {\n\t\tvar u } ${sn} {\n\t\tif err := rows.Scan(&u.ID} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            append {, &u.} ${f.pascalName()} ;
          }
        }
        append {); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\titems = append(items, u)\n\t}\n\treturn items, rows.Err()\n}\n\n} ;

        // Update
        append {func (r *} ${rn} {) Update(u *} ${sn} {) error {\n\treturn r.db.QueryRow(\n\t\t`UPDATE } ${tn} { SET } ;
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
        append {, updated_at = NOW()\n\t\t WHERE id = $} ${nextParam} {\n\t\t RETURNING updated_at`,\n} ;
        foreach field in schema.Fields {
          if (field.isInstanceOf(Field)) {
            node<Field> f = (Field) field;
            if (f.dataType != DataType.timestamp) { append {\t\tu.} ${f.pascalName()} {,\n} ; }
          }
        }
        append {\t\tu.ID,\n\t).Scan(&u.UpdatedAt)\n}\n\n} ;

        // Delete
        append {func (r *} ${rn} {) Delete(id int64) error {\n\t_, err := r.db.Exec(`DELETE FROM } ${tn} { WHERE id = $1`, id)\n\treturn err\n}\n\n} ;
      }
    }

    // ============================================================
    // Repositories — join schemas
    // ============================================================

    foreach schema in model.schemas {
      if (schema.hasReferences()) {
        string sn = schema.structName();
        string rn = schema.repoName();
        string tn = schema.name;

        append {type } ${rn} { struct{ db *sql.DB }\n\n} ;

        // Assign
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
// query
        append {) (*} ${sn} {, error) {\n\tur := &} ${sn} {{}\n\terr := r.db.QueryRow(\n\t\t`INSERT INTO } ${tn} { (} ;
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
        append {)\n\treturn ur, err\n}\n\n} ;

        // Remove
        append {func (r *} ${rn} {) Remove(} ;
        int rmIdx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (rmIdx > 0) { append {, } ; }
            append ${fr.name} { int64} ;
            rmIdx = rmIdx + 1;
          }
        }
        append {) error {\n\t_, err := r.db.Exec(`DELETE FROM } ${tn} { WHERE } ;
        int wIdx = 0;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (wIdx > 0) { append { AND } ; }
            wIdx = wIdx + 1;
            append ${fr.name} { = $} ${wIdx} ;
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

        // Cross-queries
        foreach refA in schema.Fields {
          if (refA.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> frA = (FieldRefrence) refA;
            foreach refB in schema.Fields {
              if (refB.isInstanceOf(FieldRefrence) && refB != refA) {
                node<FieldRefrence> frB = (FieldRefrence) refB;
                string ts = frB.target_schema.structName();
                string tt = frB.target_schema.name;

                append {func (r *} ${rn} {) Get} ${ts} {sBy} ${frA.target_schema.structName()} {(} ${frA.name} { int64) ([]} ${ts} {, error) {\n\trows, err := r.db.Query(\n\t\t`SELECT t.id} ;
                foreach tf in frB.target_schema.Fields {
                  if (tf.isInstanceOf(Field)) {
                    node<Field> f = (Field) tf;
                    append {, t.} ${f.name} ;
                  }
                }
                append {\n\t\t FROM } ${tt} { t\n\t\t INNER JOIN } ${tn} { j ON j.} ${frB.name} { = t.id\n\t\t WHERE j.} ${frA.name} { = $1\n\t\t ORDER BY t.id`, } ${frA.name} {,\n\t)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tdefer rows.Close()\n\tvar items []} ${ts} {\n\tfor rows.Next() {\n\t\tvar item } ${ts} {\n\t\tif err := rows.Scan(&item.ID} ;
                foreach tf in frB.target_schema.Fields {
                  if (tf.isInstanceOf(Field)) {
                    node<Field> f = (Field) tf;
                    append {, &item.} ${f.pascalName()} ;
                  }
                }
                append {); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\titems = append(items, item)\n\t}\n\treturn items, rows.Err()\n}\n\n} ;
              }
            }
          }
        }
      }
    }

    // ============================================================
    // HTTP Handlers — regular schemas
    // ============================================================

    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {
        string sn = schema.structName();
        string rn = schema.repoName();

        append {// ============================================================\n// HTTP Handlers — } ${sn} {\n// ============================================================\n\n} ;

        // Create handler
        append {func handleCreate} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\tvar u } ${sn} {\n\t\tif err := json.NewDecoder(r.Body).Decode(&u); err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tif err := repo.Create(&u); err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n\t\tw.Header().Set("Content-Type", "application/json")\n\t\tw.WriteHeader(http.StatusCreated)\n\t\tjson.NewEncoder(w).Encode(u)\n\t}\n}\n\n} ;

        // Get handler
        append {func handleGet} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\tidStr := r.PathValue("id")\n\t\tid, err := strconv.ParseInt(idStr, 10, 64)\n\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tu, err := repo.GetByID(id)\n\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusNotFound)\n\t\t\treturn\n\t\t}\n\t\tw.Header().Set("Content-Type", "application/json")\n\t\tjson.NewEncoder(w).Encode(u)\n\t}\n}\n\n} ;

        // List handler
        append {func handleList} ${sn} {s(repo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\titems, err := repo.List()\n\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n\t\tw.Header().Set("Content-Type", "application/json")\n\t\tjson.NewEncoder(w).Encode(items)\n\t}\n}\n\n} ;

        // Update handler
        append {func handleUpdate} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\tidStr := r.PathValue("id")\n\t\tid, err := strconv.ParseInt(idStr, 10, 64)\n\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tvar u } ${sn} {\n\t\tif err := json.NewDecoder(r.Body).Decode(&u); err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tu.ID = id\n\t\tif err := repo.Update(&u); err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n\t\tw.Header().Set("Content-Type", "application/json")\n\t\tjson.NewEncoder(w).Encode(u)\n\t}\n}\n\n} ;

        // Delete handler
        append {func handleDelete} ${sn} {(repo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\tidStr := r.PathValue("id")\n\t\tid, err := strconv.ParseInt(idStr, 10, 64)\n\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tif err := repo.Delete(id); err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n\t\tw.WriteHeader(http.StatusNoContent)\n\t}\n}\n\n} ;
      }
    }

    // ============================================================
    // HTTP Handlers — join schemas
    // ============================================================

    foreach schema in model.schemas {
      if (schema.hasReferences()) {
        string sn = schema.structName();
        string rn = schema.repoName();

        append {// ============================================================\n// HTTP Handlers — } ${sn} { (assignments)\n// ============================================================\n\n} ;

        node<FieldRefrence> firstRef = null;
        node<FieldRefrence> secondRef = null;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (firstRef == null) { firstRef = fr; }
            else if (secondRef == null) { secondRef = fr; }
          }
        }

        // Assign handler
        append {func handleAssign} ${secondRef.target_schema.structName()} {(urRepo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\t} ${firstRef.name} {, err := strconv.ParseInt(r.PathValue("id"), 10, 64)\n\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tvar body Assign} ${secondRef.target_schema.structName()} {Body\n\t\tif err := json.NewDecoder(r.Body).Decode(&body); err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tur, err := urRepo.Assign(} ${firstRef.name} {, body.} ${secondRef.pascalName()} {)\n\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n\t\tw.Header().Set("Content-Type", "application/json")\n\t\tw.WriteHeader(http.StatusCreated)\n\t\tjson.NewEncoder(w).Encode(ur)\n\t}\n}\n\n} ;

        // Remove handler
        append {func handleRemove} ${secondRef.target_schema.structName()} {(urRepo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\t} ${firstRef.name} {, err := strconv.ParseInt(r.PathValue("id"), 10, 64)\n\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\t} ${secondRef.name} {, err := strconv.ParseInt(r.PathValue("} ${secondRef.name} {"), 10, 64)\n\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\tif err := urRepo.Remove(} ${firstRef.name} {, } ${secondRef.name} {); err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n\t\tw.WriteHeader(http.StatusNoContent)\n\t}\n}\n\n} ;

        // Cross-query handlers
        foreach refA in schema.Fields {
          if (refA.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> frA = (FieldRefrence) refA;
            foreach refB in schema.Fields {
              if (refB.isInstanceOf(FieldRefrence) && refB != refA) {
                node<FieldRefrence> frB = (FieldRefrence) refB;
                append {func handleGet} ${frA.target_schema.structName()} ${frB.target_schema.structName()} {s(urRepo *} ${rn} {) http.HandlerFunc {\n\treturn func(w http.ResponseWriter, r *http.Request) {\n\t\tid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)\n\t\tif err != nil {\n\t\t\thttp.Error(w, "invalid id", http.StatusBadRequest)\n\t\t\treturn\n\t\t}\n\t\titems, err := urRepo.Get} ${frB.target_schema.structName()} {sBy} ${frA.target_schema.structName()} {(id)\n\t\tif err != nil {\n\t\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n\t\t\treturn\n\t\t}\n\t\tw.Header().Set("Content-Type", "application/json")\n\t\tjson.NewEncoder(w).Encode(items)\n\t}\n}\n\n} ;
              }
            }
          }
        }
      }
    }

    // ============================================================
    // Main
    // ============================================================

    append {// ============================================================\n// Main\n// ============================================================\n\nfunc main() {\n\tdbURL := os.Getenv("DATABASE_URL")\n\tif dbURL == "" {\n\t\tdbURL = "postgres://} ${infra.dbUser} {:} ${infra.dbPass} {@localhost:5432/} ${infra.dbName} {?sslmode=disable"\n\t}\n\n\tdb, err := sql.Open("postgres", dbURL)\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\tdefer db.Close()\n\n\tfor i := 0; i < 5; i++ {\n\t\tif err = db.Ping(); err == nil {\n\t\t\tbreak\n\t\t}\n\t\tlog.Printf("DB not ready, retrying... (%d/5)", i+1)\n\t\ttime.Sleep(2 * time.Second)\n\t}\n\tif err != nil {\n\t\tlog.Fatal("DB connection failed:", err)\n\t}\n\n\tif _, err := db.Exec(migrationSQL); err != nil {\n\t\tlog.Fatal(err)\n\t}\n\tlog.Println("Migration complete")\n\n} ;

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

    append {\n\tmux := http.NewServeMux()\n\n\t// Swagger UI\n\tmux.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)\n\n} ;

    // Regular routes
    foreach schema in model.schemas {
      if (!(schema.hasReferences())) {
        string vr = schema.singularName() + "Repo";
        string sn = schema.structName();
        string tn = schema.name;
        append {\t// } ${sn} {s\n\tmux.HandleFunc("POST /} ${tn} {", handleCreate} ${sn} {(} ${vr} {))\n\tmux.HandleFunc("GET /} ${tn} {", handleList} ${sn} {s(} ${vr} {))\n\tmux.HandleFunc("GET /} ${tn} {/{id}", handleGet} ${sn} {(} ${vr} {))\n\tmux.HandleFunc("PUT /} ${tn} {/{id}", handleUpdate} ${sn} {(} ${vr} {))\n\tmux.HandleFunc("DELETE /} ${tn} {/{id}", handleDelete} ${sn} {(} ${vr} {))\n\n} ;
      }
    }

    // Join routes
    foreach schema in model.schemas {
      if (schema.hasReferences()) {
        string vr = schema.singularName() + "Repo";
        node<FieldRefrence> fRef = null;
        node<FieldRefrence> sRef = null;
        foreach field in schema.Fields {
          if (field.isInstanceOf(FieldRefrence)) {
            node<FieldRefrence> fr = (FieldRefrence) field;
            if (fRef == null) { fRef = fr; }
            else if (sRef == null) { sRef = fr; }
          }
        }
        append {\t// } ${schema.structName()} { assignments\n\tmux.HandleFunc("POST /} ${fRef.target_schema.name} {/{id}/} ${sRef.target_schema.name} {", handleAssign} ${sRef.target_schema.structName()} {(} ${vr} {))\n\tmux.HandleFunc("DELETE /} ${fRef.target_schema.name} {/{id}/} ${sRef.target_schema.name} {/{} ${sRef.name} {}", handleRemove} ${sRef.target_schema.structName()} {(} ${vr} {))\n\tmux.HandleFunc("GET /} ${fRef.target_schema.name} {/{id}/} ${sRef.target_schema.name} {", handleGet} ${fRef.target_schema.structName()} ${sRef.target_schema.structName()} {s(} ${vr} {))\n\tmux.HandleFunc("GET /} ${sRef.target_schema.name} {/{id}/} ${fRef.target_schema.name} {", handleGet} ${sRef.target_schema.structName()} ${fRef.target_schema.structName()} {s(} ${vr} {))\n\n} ;
      }
    }

    append {\tfmt.Println("Serving on :} ${infra.port} {")\n\tfmt.Println("Swagger UI: http://localhost:} ${infra.port} {/swagger/index.html")\n\tlog.Fatal(http.ListenAndServe(":} ${infra.port} {", mux))\n}\n} ;

  }
  file path : (node)->string { "go/"; }
  extension : (node)->string { "go"; }
}
```
