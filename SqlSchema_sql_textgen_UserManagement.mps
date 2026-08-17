<?xml version="1.0" encoding="UTF-8"?>
<model ref="r:f10fb75e-9e61-4ecf-8328-e03d2ea8b365(UserManagement.textGen)">
  <persistence version="9" />
  <languages>
    <use id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen" version="1" />
    <devkit ref="fa73d85a-ac7f-447b-846c-fcdc41caa600(jetbrains.mps.devkit.aspect.textgen)" />
  </languages>
  <imports>
    <import index="s1" ref="r:a3a366a2-da30-48fe-b644-04a6d92b06a4(UserManagement.structure)" implicit="true" />
  </imports>
  <registry>
    <language id="f3061a53-9226-4cc5-a443-f952ceaf5816" name="jetbrains.mps.baseLanguage">
      <concept id="1137021947720" name="jetbrains.mps.baseLanguage.structure.ConceptFunction" flags="in" index="2VMwT0">
        <child id="1137022507850" name="body" index="2VODD2" />
      </concept>
      <concept id="1070475926800" name="jetbrains.mps.baseLanguage.structure.StringLiteral" flags="nn" index="Xl_RD">
        <property id="1070475926801" name="value" index="Xl_RC" />
      </concept>
      <concept id="1068580123157" name="jetbrains.mps.baseLanguage.structure.Statement" flags="nn" index="3clFbH" />
      <concept id="1068580123136" name="jetbrains.mps.baseLanguage.structure.StatementList" flags="sn" stub="5293379017992965193" index="3clFbS">
        <child id="1068581517665" name="statement" index="3cqZAp" />
      </concept>
      <concept id="1068581242878" name="jetbrains.mps.baseLanguage.structure.ReturnStatement" flags="nn" index="3cpWs6">
        <child id="1068581517676" name="expression" index="3cqZAk" />
      </concept>
    </language>
    <language id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen">
      <concept id="45307784116571022" name="jetbrains.mps.lang.textGen.structure.FilenameFunction" flags="ig" index="29tfMY" />
      <concept id="1237305208784" name="jetbrains.mps.lang.textGen.structure.NewLineAppendPart" flags="ng" index="l8MVK" />
      <concept id="1237305557638" name="jetbrains.mps.lang.textGen.structure.ConstantStringAppendPart" flags="ng" index="la8eA">
        <property id="1237305576108" name="value" index="lacIc" />
      </concept>
      <concept id="1237306079178" name="jetbrains.mps.lang.textGen.structure.AppendOperation" flags="nn" index="lc7rE">
        <child id="1237306115446" name="part" index="lcghm" />
      </concept>
      <concept id="1233670071145" name="jetbrains.mps.lang.textGen.structure.ConceptTextGenDeclaration" flags="ig" index="WtQ9Q">
        <reference id="1233670257997" name="conceptDeclaration" index="WuzLi" />
        <child id="45307784116711884" name="filename" index="29tGrW" />
        <child id="1233749296504" name="textGenBlock" index="11c4hB" />
      </concept>
      <concept id="1233748055915" name="jetbrains.mps.lang.textGen.structure.NodeParameter" flags="nn" index="117lpO" />
      <concept id="1233749247888" name="jetbrains.mps.lang.textGen.structure.GenerateTextDeclaration" flags="in" index="11bSqf" />
    </language>
  </registry>
  <node concept="WtQ9Q" id="7tgPrsAeQ">
    <ref role="WuzLi" to="s1:7669448123830914766" resolve="SqlSchem" />
    <node concept="29tfMY" id="7tgPrsAeT" role="29tGrW">
      <node concept="3clFbS" id="7tgPrsAeU" role="2VODD2">
        <node concept="3cpWs6" id="7tgPrsAeV" role="3cqZAp">
          <node concept="Xl_RD" id="7tgPrsAeW" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="7tgPrsAeR" role="11c4hB">
      <node concept="3clFbS" id="7tgPrsAeS" role="2VODD2">
        <node concept="3clFbH" id="7tgPrsAa5" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAa8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa6" role="lcghm">
            <property role="lacIc" value="-- Auto-generated SQL schema" />
          </node>
          <node concept="l8MVK" id="7tgPrsAa7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa9" role="lcghm">
            <property role="lacIc" value="-- Schema: " />
          </node>
          <node concept="la8eA" id="7tgPrsAba" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbe" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbf" role="lcghm">
            <property role="lacIc" value="CREATE SCHEMA IF NOT EXISTS " />
          </node>
          <node concept="la8eA" id="7tgPrsAbg" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbh" role="lcghm">
            <property role="lacIc" value=";" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbi" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbk" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAbm" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAbn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbo" role="lcghm">
            <property role="lacIc" value="{???-foreach ref in node.entityrefs {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbq" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Entity&gt; entity = ref.entity;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAbr" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAbx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbs" role="lcghm">
            <property role="lacIc" value="-- Table: " />
          </node>
          <node concept="la8eA" id="7tgPrsAbt" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbu" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAbv" role="lcghm">
            <property role="lacIc" value="{???-entity.tableName}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAby" role="lcghm">
            <property role="lacIc" value="CREATE TABLE IF NOT EXISTS " />
          </node>
          <node concept="la8eA" id="7tgPrsAbz" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbA" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAbB" role="lcghm">
            <property role="lacIc" value="{???-entity.tableName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbC" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbD" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAbF" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAbG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbH" role="lcghm">
            <property role="lacIc" value="{???-int fieldIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbJ" role="lcghm">
            <property role="lacIc" value="{???-foreach field in entity.fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbL" role="lcghm">
            <property role="lacIc" value="{???-if (fieldIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbM" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAbN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbQ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbS" role="lcghm">
            <property role="lacIc" value="{???-if (field.hasAnnotation(FieldAnnotation:primaryKey)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbT" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAbU" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbV" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAbW" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbX" role="lcghm">
            <property role="lacIc" value=" PRIMARY KEY DEFAULT gen_random_uuid()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb0" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAb1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb2" role="lcghm">
            <property role="lacIc" value="{???-if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:nullable)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:auto)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:unique))) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAb8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb3" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAb4" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb5" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAb6" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb7" role="lcghm">
            <property role="lacIc" value=" NOT NULL" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAb9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAca" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcc" role="lcghm">
            <property role="lacIc" value="{???-if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:nullable)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:auto)) &amp;&amp; field.hasAnnotation(FieldAnnotation:unique)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAci" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcd" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAce" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcf" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAcg" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAch" role="lcghm">
            <property role="lacIc" value=" NOT NULL UNIQUE" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAck" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcm" role="lcghm">
            <property role="lacIc" value="{???-if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) &amp;&amp; field.hasAnnotation(FieldAnnotation:nullable)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcn" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAco" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcp" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAcq" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAct" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcv" role="lcghm">
            <property role="lacIc" value="{???-if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) &amp;&amp; field.hasAnnotation(FieldAnnotation:auto)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcw" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcx" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcy" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAcz" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcA" role="lcghm">
            <property role="lacIc" value=" NOT NULL DEFAULT NOW()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcD" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcF" role="lcghm">
            <property role="lacIc" value="{???-fieldIdx = fieldIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcH" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAcI" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcL" role="lcghm">
            <property role="lacIc" value=");" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcO" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAcQ" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcS" role="lcghm">
            <property role="lacIc" value="{???-foreach field in entity.fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcU" role="lcghm">
            <property role="lacIc" value="{???-if (field.hasAnnotation(FieldAnnotation:indexed)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcV" role="lcghm">
            <property role="lacIc" value="CREATE INDEX IF NOT EXISTS idx_" />
          </node>
          <node concept="la8eA" id="7tgPrsAcW" role="lcghm">
            <property role="lacIc" value="{???-entity.tableName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcX" role="lcghm">
            <property role="lacIc" value="_" />
          </node>
          <node concept="la8eA" id="7tgPrsAcY" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcZ" role="lcghm">
            <property role="lacIc" value=" ON " />
          </node>
          <node concept="la8eA" id="7tgPrsAc0" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAc1" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAc2" role="lcghm">
            <property role="lacIc" value="{???-entity.tableName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAc3" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAc4" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAc5" role="lcghm">
            <property role="lacIc" value=");" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc9" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAda" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdb" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdc" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAde" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdg" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdh" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdj" role="lcghm">
            <property role="lacIc" value="{???-foreach relation in node.relations {}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdk" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdl" role="lcghm">
            <property role="lacIc" value="-- Junction table: " />
          </node>
          <node concept="la8eA" id="7tgPrsAdm" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdn" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAdo" role="lcghm">
            <property role="lacIc" value="{???-relation.tableName}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdr" role="lcghm">
            <property role="lacIc" value="CREATE TABLE IF NOT EXISTS " />
          </node>
          <node concept="la8eA" id="7tgPrsAds" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdt" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAdu" role="lcghm">
            <property role="lacIc" value="{???-relation.tableName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdv" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdy" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdz" role="lcghm">
            <property role="lacIc" value="{???-relation.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdA" role="lcghm">
            <property role="lacIc" value="_id UUID NOT NULL REFERENCES " />
          </node>
          <node concept="la8eA" id="7tgPrsAdB" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdC" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAdD" role="lcghm">
            <property role="lacIc" value="{???-relation.from.tableName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdE" role="lcghm">
            <property role="lacIc" value="(id) ON DELETE CASCADE," />
          </node>
          <node concept="l8MVK" id="7tgPrsAdF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdH" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdI" role="lcghm">
            <property role="lacIc" value="{???-relation.to.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdJ" role="lcghm">
            <property role="lacIc" value="_id UUID NOT NULL REFERENCES " />
          </node>
          <node concept="la8eA" id="7tgPrsAdK" role="lcghm">
            <property role="lacIc" value="{???-node.dbSchema}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdL" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAdM" role="lcghm">
            <property role="lacIc" value="{???-relation.to.tableName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdN" role="lcghm">
            <property role="lacIc" value="(id) ON DELETE CASCADE" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdP" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdR" role="lcghm">
            <property role="lacIc" value="{???-foreach field in relation.extraFields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdT" role="lcghm">
            <property role="lacIc" value="{???-if (field.hasAnnotation(FieldAnnotation:auto)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdU" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAdV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAdW" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdX" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdY" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAdZ" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAd0" role="lcghm">
            <property role="lacIc" value=" NOT NULL DEFAULT NOW()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd3" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd5" role="lcghm">
            <property role="lacIc" value="{???-if (!(field.hasAnnotation(FieldAnnotation:auto)) &amp;&amp; field.hasAnnotation(FieldAnnotation:nullable)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAec" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd6" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAd7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAd8" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAd9" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAea" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAeb" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAed" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAee" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAef" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeg" role="lcghm">
            <property role="lacIc" value="{???-if (!(field.hasAnnotation(FieldAnnotation:auto)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:nullable))) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeh" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAei" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAej" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAek" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAel" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAem" role="lcghm">
            <property role="lacIc" value="{???-field.sqlType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAen" role="lcghm">
            <property role="lacIc" value=" NOT NULL" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAep" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeq" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAer" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAes" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAet" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAeC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeu" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAev" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAew" role="lcghm">
            <property role="lacIc" value="	PRIMARY KEY (" />
          </node>
          <node concept="la8eA" id="7tgPrsAex" role="lcghm">
            <property role="lacIc" value="{???-relation.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAey" role="lcghm">
            <property role="lacIc" value="_id, " />
          </node>
          <node concept="la8eA" id="7tgPrsAez" role="lcghm">
            <property role="lacIc" value="{???-relation.to.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeA" role="lcghm">
            <property role="lacIc" value="_id)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeD" role="lcghm">
            <property role="lacIc" value=");" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeH" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeG" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAeI" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAeJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeK" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAeL" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAeM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeN" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeP" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
      </node>
    </node>
  </node>
</model>
