```
text gen component for concept NatsServer {
file name :
  init
file path :
extension :
  sql
(node)->void {


// ─── header ───
append {-- Auto-generated SQL for } ${node.name} {\n} ;
append {-- Schema: } ${node.dbSchema} {\n} ;
append {\n} ;
append {CREATE SCHEMA IF NOT EXISTS } ${node.dbSchema} {;\n} ;
append {\n} ;

// ─── Entity tables ───
foreach entity in node.entities {
append {-- Table: } ${node.dbSchema} {.} ${entity.tableName} {\n} ;
append {CREATE TABLE IF NOT EXISTS } ${node.dbSchema} {.} ${entity.tableName} { (\n} ;

int fieldIdx = 0;
foreach field in entity.fields {

if (fieldIdx > 0) {
append {,\n} ;
}

if (field.hasAnnotation(FieldAnnotation:primaryKey)) {
append {\t} ${field.dbName()} { } ${field.sqlType()} { PRIMARY KEY DEFAULT gen_random_uuid()} ;
}
if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) && !(field.hasAnnotation(FieldAnnotation:nullable)) && !(field.hasAnnotation(FieldAnnotation:auto)) && !(field.hasAnnotation(FieldAnnotation:unique))) {
append {\t} ${field.dbName()} { } ${field.sqlType()} { NOT NULL} ;
}
if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) && !(field.hasAnnotation(FieldAnnotation:nullable)) && !(field.hasAnnotation(FieldAnnotation:auto)) && field.hasAnnotation(FieldAnnotation:unique)) {
append {\t} ${field.dbName()} { } ${field.sqlType()} { NOT NULL UNIQUE} ;
}
if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) && field.hasAnnotation(FieldAnnotation:nullable)) {
append {\t} ${field.dbName()} { } ${field.sqlType()} ;
}
if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) && field.hasAnnotation(FieldAnnotation:auto)) {
append {\t} ${field.dbName()} { } ${field.sqlType()} { NOT NULL DEFAULT NOW()} ;
}

fieldIdx = fieldIdx + 1;
}

append {\n} ;
append {);\n} ;
append {\n} ;

// ─── Indexes for indexed fields ───
foreach field in entity.fields {
if (field.hasAnnotation(FieldAnnotation:indexed)) {
append {CREATE INDEX IF NOT EXISTS idx_} ${entity.tableName} {_} ${field.dbName()} { ON } ${node.dbSchema} {.} ${entity.tableName} {(} ${field.dbName()} {);\n} ;
}
}
append {\n} ;
}

// ─── Relation junction tables ───
foreach relation in node.relations {
append {-- Junction table: } ${node.dbSchema} {.} ${relation.tableName} {\n} ;
append {CREATE TABLE IF NOT EXISTS } ${node.dbSchema} {.} ${relation.tableName} { (\n} ;
append {\t} ${relation.from.name.toLowerCase()} {_id UUID NOT NULL REFERENCES } ${node.dbSchema} {.} ${relation.from.tableName} {(id) ON DELETE CASCADE,\n} ;
append {\t} ${relation.to.name.toLowerCase()} {_id UUID NOT NULL REFERENCES } ${node.dbSchema} {.} ${relation.to.tableName} {(id) ON DELETE CASCADE} ;

foreach field in relation.extraFields {
append {,\n} ;
if (field.hasAnnotation(FieldAnnotation:auto)) {
append {\t} ${field.dbName()} { } ${field.sqlType()} { NOT NULL DEFAULT NOW()} ;
}
if (field.hasAnnotation(FieldAnnotation:nullable)) {
append {\t} ${field.dbName()} { } ${field.sqlType()} ;
}
if (!(field.hasAnnotation(FieldAnnotation:auto)) && !(field.hasAnnotation(FieldAnnotation:nullable))) {
append {\t} ${field.dbName()} { } ${field.sqlType()} { NOT NULL} ;
}
}

append {,\n} ;
append {\tPRIMARY KEY (} ${relation.from.name.toLowerCase()} {_id, } ${relation.to.name.toLowerCase()} {_id)\n} ;
append {);\n} ;
append {\n} ;
}
```

}
}
