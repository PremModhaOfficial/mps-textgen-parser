<?xml version="1.0" encoding="UTF-8"?>
<model ref="TODO-MODEL-UUID">
  <persistence version="9" />
  <languages>
    <use id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen" version="1" />
    <devkit ref="fa73d85a-ac7f-447b-846c-fcdc41caa600(jetbrains.mps.devkit.aspect.textgen)" />
  </languages>
  <imports>
    <import index="s1" ref="TODO-STRUCTURE-UUID" implicit="true" />
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
  <node concept="WtQ9Q" id="7tgPrsAES">
    <ref role="WuzLi" to="s1:TODO-CONCEPT-ID" resolve="Code" />
    <node concept="29tfMY" id="7tgPrsAEV" role="29tGrW">
      <node concept="3clFbS" id="7tgPrsAEW" role="2VODD2">
        <node concept="3cpWs6" id="7tgPrsAEX" role="3cqZAp">
          <node concept="Xl_RD" id="7tgPrsAEY" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="7tgPrsAET" role="11c4hB">
      <node concept="3clFbS" id="7tgPrsAEU" role="2VODD2">
        <node concept="lc7rE" id="7tgPrsAa5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa6" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Code&gt; n = node;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAa7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa8" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Models&gt; model = n.model;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAa9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAba" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Infra&gt; infra = n.infra;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAbb" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAci" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbc" role="lcghm">
            <property role="lacIc" value="package main" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbd" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbe" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbf" role="lcghm">
            <property role="lacIc" value="import (" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbh" role="lcghm">
            <property role="lacIc" value="	&quot;database/sql&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbi" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbj" role="lcghm">
            <property role="lacIc" value="	_ &quot;embed&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbl" role="lcghm">
            <property role="lacIc" value="	&quot;encoding/json&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbn" role="lcghm">
            <property role="lacIc" value="	&quot;fmt&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbo" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbp" role="lcghm">
            <property role="lacIc" value="	&quot;log&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbr" role="lcghm">
            <property role="lacIc" value="	&quot;net/http&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbs" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbt" role="lcghm">
            <property role="lacIc" value="	&quot;os&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbu" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbv" role="lcghm">
            <property role="lacIc" value="	&quot;strconv&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbw" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbx" role="lcghm">
            <property role="lacIc" value="	&quot;time&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAby" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbA" role="lcghm">
            <property role="lacIc" value="	_ &quot;github.com/lib/pq&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbC" role="lcghm">
            <property role="lacIc" value="	httpSwagger &quot;github.com/swaggo/http-swagger&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbE" role="lcghm">
            <property role="lacIc" value="	_ &quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAbF" role="lcghm">
            <property role="lacIc" value="{???-infra.modulePath}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbG" role="lcghm">
            <property role="lacIc" value="/docs&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbI" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbJ" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbL" role="lcghm">
            <property role="lacIc" value="//go:embed user_management_init.sql" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbN" role="lcghm">
            <property role="lacIc" value="var migrationSQL string" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbO" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbQ" role="lcghm">
            <property role="lacIc" value="// @title         " />
          </node>
          <node concept="la8eA" id="7tgPrsAbR" role="lcghm">
            <property role="lacIc" value="{???-model.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbS" role="lcghm">
            <property role="lacIc" value=" API" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbU" role="lcghm">
            <property role="lacIc" value="// @version       1.0" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbW" role="lcghm">
            <property role="lacIc" value="// @description   CRUD service for " />
          </node>
          <node concept="la8eA" id="7tgPrsAbX" role="lcghm">
            <property role="lacIc" value="{???-model.name}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbZ" role="lcghm">
            <property role="lacIc" value="// @host          localhost:" />
          </node>
          <node concept="la8eA" id="7tgPrsAb0" role="lcghm">
            <property role="lacIc" value="{???-infra.port}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb2" role="lcghm">
            <property role="lacIc" value="// @BasePath      /" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb4" role="lcghm">
            <property role="lacIc" value="// @schemes       http" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb6" role="lcghm">
            <property role="lacIc" value="// @produce       json" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb8" role="lcghm">
            <property role="lacIc" value="// @consumes      json" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb9" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAca" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcb" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcd" role="lcghm">
            <property role="lacIc" value="// Models" />
          </node>
          <node concept="l8MVK" id="7tgPrsAce" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcf" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcg" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAch" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAcj" role="3cqZAp" />
        <!-- // ~~~~ Regular schema structs ~~~~ -->
        <node concept="lc7rE" id="7tgPrsAck" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcl" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcn" role="lcghm">
            <property role="lacIc" value="{???-if (!(schema.hasReferences())) {}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAco" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcp" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAcq" role="lcghm">
            <property role="lacIc" value="{???-schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcr" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcs" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAct" role="lcghm">
            <property role="lacIc" value="	ID int64 `json:&quot;id&quot; db:&quot;id&quot; example:&quot;1&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcx" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcz" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcB" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcC" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcD" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcE" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAcF" role="lcghm">
            <property role="lacIc" value="{???-f.goType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcG" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcH" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcI" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcJ" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcK" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcO" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcQ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcR" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcS" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAcT" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAcV" role="3cqZAp" />
        <!-- // MarshalJSON if Sensitive -->
        <node concept="lc7rE" id="7tgPrsAcW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcX" role="lcghm">
            <property role="lacIc" value="{???-boolean hasSensitive = false;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcZ" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc1" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc3" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc5" role="lcghm">
            <property role="lacIc" value="{???-if (f.Sensitive) { hasSensitive = true; } }" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc7" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc9" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAda" role="3cqZAp" />
        <!-- // found the sensitive parst -->
        <node concept="lc7rE" id="7tgPrsAdb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdc" role="lcghm">
            <property role="lacIc" value="{???-if (hasSensitive) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdd" role="lcghm">
            <property role="lacIc" value="func (u " />
          </node>
          <node concept="la8eA" id="7tgPrsAde" role="lcghm">
            <property role="lacIc" value="{???-schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdf" role="lcghm">
            <property role="lacIc" value=") MarshalJSON() ([]byte, error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAdh" role="lcghm">
            <property role="lacIc" value="	type Alias " />
          </node>
          <node concept="la8eA" id="7tgPrsAdi" role="lcghm">
            <property role="lacIc" value="{???-schema.structName()}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAdk" role="lcghm">
            <property role="lacIc" value="	return json.Marshal(&amp;struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdl" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAdm" role="lcghm">
            <property role="lacIc" value="		Alias" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdn" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdq" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAds" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdu" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdw" role="lcghm">
            <property role="lacIc" value="{???-if (f.Sensitive) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdx" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAdy" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdz" role="lcghm">
            <property role="lacIc" value=" string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdA" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdB" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdF" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdH" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdJ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdK" role="lcghm">
            <property role="lacIc" value="	}{" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAdM" role="lcghm">
            <property role="lacIc" value="		Alias: (Alias)(u)," />
          </node>
          <node concept="l8MVK" id="7tgPrsAdN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdQ" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdS" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdU" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdW" role="lcghm">
            <property role="lacIc" value="{???-if (f.Sensitive) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdX" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAdY" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdZ" role="lcghm">
            <property role="lacIc" value=": &quot;[REDACTED]&quot;," />
          </node>
          <node concept="l8MVK" id="7tgPrsAd0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd3" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd5" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd7" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAed" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd8" role="lcghm">
            <property role="lacIc" value="	})" />
          </node>
          <node concept="l8MVK" id="7tgPrsAd9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAea" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeb" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAec" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAee" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAef" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAeg" role="3cqZAp" />
        <node concept="3clFbH" id="7tgPrsAeh" role="3cqZAp" />
        <!-- // Create struct -->
        <node concept="lc7rE" id="7tgPrsAem" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAei" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAej" role="lcghm">
            <property role="lacIc" value="{???-schema.createStructName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAek" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAel" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAen" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeo" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAep" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeq" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAer" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAes" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAet" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeu" role="lcghm">
            <property role="lacIc" value="{???-if (f.dataType != DataType.timestamp) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAev" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAew" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAex" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAey" role="lcghm">
            <property role="lacIc" value="{???-f.goType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAez" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAeA" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeB" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeF" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeH" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeJ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeK" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeL" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAeM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeP" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeR" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAeS" role="3cqZAp" />
        <!-- // ~~~~ Join schema structs ~~~~ -->
        <node concept="lc7rE" id="7tgPrsAeT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeU" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeW" role="lcghm">
            <property role="lacIc" value="{???-if (schema.hasReferences()) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeX" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAeY" role="lcghm">
            <property role="lacIc" value="{???-schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeZ" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAe2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe3" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe5" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe7" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe8" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAe9" role="lcghm">
            <property role="lacIc" value="{???-fr.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfa" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAfb" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfc" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAfd" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfe" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAff" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfi" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfk" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfm" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfn" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAfo" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfp" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAfq" role="lcghm">
            <property role="lacIc" value="{???-f.goType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfr" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAfs" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAft" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAfu" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfv" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfz" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfB" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfC" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfD" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAfE" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfG" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAfH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfI" role="lcghm">
            <property role="lacIc" value="{???-int refCount = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfK" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfM" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfO" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfQ" role="lcghm">
            <property role="lacIc" value="{???-if (refCount == 1) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfR" role="lcghm">
            <property role="lacIc" value="type Assign" />
          </node>
          <node concept="la8eA" id="7tgPrsAfS" role="lcghm">
            <property role="lacIc" value="{???-fr.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfT" role="lcghm">
            <property role="lacIc" value="Body struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfU" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfV" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAfW" role="lcghm">
            <property role="lacIc" value="{???-fr.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfX" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAfY" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfZ" role="lcghm">
            <property role="lacIc" value="&quot; binding:&quot;required&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAf1" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf2" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAf3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAf5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf6" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf8" role="lcghm">
            <property role="lacIc" value="{???-refCount = refCount + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAga" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgc" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAge" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgg" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAgh" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — regular -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAgi" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAgq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgj" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgl" role="lcghm">
            <property role="lacIc" value="// Repositories" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgn" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgo" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAgp" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgr" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAgs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgt" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgv" role="lcghm">
            <property role="lacIc" value="{???-if (!(schema.hasReferences())) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgx" role="lcghm">
            <property role="lacIc" value="{???-string sn = schema.structName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgz" role="lcghm">
            <property role="lacIc" value="{???-string rn = schema.repoName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgB" role="lcghm">
            <property role="lacIc" value="{???-string tn = schema.name;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAgC" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAgI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgD" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAgE" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgF" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgG" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAgH" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgJ" role="3cqZAp" />
        <!-- // Create -->
        <node concept="lc7rE" id="7tgPrsAgV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgK" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAgL" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgM" role="lcghm">
            <property role="lacIc" value=") Create(u *" />
          </node>
          <node concept="la8eA" id="7tgPrsAgN" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgO" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgQ" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgS" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
          <node concept="la8eA" id="7tgPrsAgT" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgU" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgX" role="lcghm">
            <property role="lacIc" value="{???-int idx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgZ" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg1" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg3" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg5" role="lcghm">
            <property role="lacIc" value="{???-if (f.dataType != DataType.timestamp) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg7" role="lcghm">
            <property role="lacIc" value="{???-if (idx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg8" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAha" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhb" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhc" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhf" role="lcghm">
            <property role="lacIc" value="{???-idx = idx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhh" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhj" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhl" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhm" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAho" role="lcghm">
            <property role="lacIc" value="		 VALUES (" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhr" role="lcghm">
            <property role="lacIc" value="{???-for (int i = 1; i &lt;= idx; i++) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAht" role="lcghm">
            <property role="lacIc" value="{???-if (i &gt; 1) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhu" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhx" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhy" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
          <node concept="la8eA" id="7tgPrsAhz" role="lcghm">
            <property role="lacIc" value="{???-i}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhC" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhD" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAhF" role="lcghm">
            <property role="lacIc" value="		 RETURNING id" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhI" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhK" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhM" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhO" role="lcghm">
            <property role="lacIc" value="{???-if (f.dataType == DataType.timestamp) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhP" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAhQ" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhT" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhV" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhX" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAh0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhY" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
          <node concept="l8MVK" id="7tgPrsAhZ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAh1" role="3cqZAp" />
        <!-- // non timestapm -->
        <node concept="lc7rE" id="7tgPrsAh2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh3" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAh4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh5" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAh6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh7" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAh8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh9" role="lcghm">
            <property role="lacIc" value="{???-if (f.dataType != DataType.timestamp) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAie" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAia" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
          <node concept="la8eA" id="7tgPrsAib" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAic" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAid" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAif" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAig" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAih" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAii" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAij" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAik" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAil" role="3cqZAp" />
        <!-- // scanning -->
        <node concept="lc7rE" id="7tgPrsAin" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAim" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAio" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAip" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAir" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAis" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAit" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiv" role="lcghm">
            <property role="lacIc" value="{???-if (f.dataType == DataType.timestamp) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiw" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
          <node concept="la8eA" id="7tgPrsAix" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiA" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiC" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiE" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiF" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAiH" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiI" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAiJ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAiL" role="3cqZAp" />
        <!-- // GetByID -->
        <node concept="lc7rE" id="7tgPrsAiZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiM" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAiN" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAiO" role="lcghm">
            <property role="lacIc" value=") GetByID(id int64) (*" />
          </node>
          <node concept="la8eA" id="7tgPrsAiP" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAiQ" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAiS" role="lcghm">
            <property role="lacIc" value="	u := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAiT" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAiU" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAiW" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAiY" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi1" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi3" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi5" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi6" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAi7" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAja" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjc" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAje" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAjf" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjg" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id," />
          </node>
          <node concept="l8MVK" id="7tgPrsAjh" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAji" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjl" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjn" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjp" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjq" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
          <node concept="la8eA" id="7tgPrsAjr" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAju" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjw" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjx" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjy" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjz" role="lcghm">
            <property role="lacIc" value="	return u, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjA" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjB" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjC" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAjD" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAjF" role="3cqZAp" />
        <!-- // List -->
        <node concept="lc7rE" id="7tgPrsAjP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjG" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAjH" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjI" role="lcghm">
            <property role="lacIc" value=") List() ([]" />
          </node>
          <node concept="la8eA" id="7tgPrsAjJ" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjK" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjM" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjO" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjR" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjT" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjV" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjW" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAjX" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj0" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAj1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj2" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkp" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAj3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAj4" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAj5" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAj6" role="lcghm">
            <property role="lacIc" value=" ORDER BY id`)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAj7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAj8" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAj9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAka" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkc" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAke" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkf" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkg" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
          <node concept="la8eA" id="7tgPrsAkh" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAki" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkj" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkl" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
          <node concept="la8eA" id="7tgPrsAkm" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAko" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;u.ID" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkr" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAks" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkt" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAku" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkv" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAky" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkw" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
          <node concept="la8eA" id="7tgPrsAkx" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkA" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkC" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkD" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkF" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkH" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkJ" role="lcghm">
            <property role="lacIc" value="		items = append(items, u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkL" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkN" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkO" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkP" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkQ" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAkR" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAkT" role="3cqZAp" />
        <!-- // Update -->
        <node concept="lc7rE" id="7tgPrsAk5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkU" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAkV" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAkW" role="lcghm">
            <property role="lacIc" value=") Update(u *" />
          </node>
          <node concept="la8eA" id="7tgPrsAkX" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAkY" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAk0" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAk1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAk2" role="lcghm">
            <property role="lacIc" value="		`UPDATE " />
          </node>
          <node concept="la8eA" id="7tgPrsAk3" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAk4" role="lcghm">
            <property role="lacIc" value=" SET " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAk6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk7" role="lcghm">
            <property role="lacIc" value="{???-int uidx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAk8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk9" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAla" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlb" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAld" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAle" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlf" role="lcghm">
            <property role="lacIc" value="{???-if (f.dataType != DataType.timestamp) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlh" role="lcghm">
            <property role="lacIc" value="{???-if (uidx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAli" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAll" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAln" role="lcghm">
            <property role="lacIc" value="{???-uidx = uidx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlo" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAlp" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
          <node concept="la8eA" id="7tgPrsAlq" role="lcghm">
            <property role="lacIc" value="{???-uidx}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAls" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlt" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlv" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlx" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAly" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlz" role="lcghm">
            <property role="lacIc" value="{???-int nextParam = uidx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlA" role="lcghm">
            <property role="lacIc" value=", updated_at = NOW()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlC" role="lcghm">
            <property role="lacIc" value="		 WHERE id = $" />
          </node>
          <node concept="la8eA" id="7tgPrsAlD" role="lcghm">
            <property role="lacIc" value="{???-nextParam}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlF" role="lcghm">
            <property role="lacIc" value="		 RETURNING updated_at`," />
          </node>
          <node concept="l8MVK" id="7tgPrsAlG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlJ" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlL" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlN" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlP" role="lcghm">
            <property role="lacIc" value="{???-if (f.dataType != DataType.timestamp) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlQ" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
          <node concept="la8eA" id="7tgPrsAlR" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAlS" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAlT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlW" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlY" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAl0" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAl8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAl1" role="lcghm">
            <property role="lacIc" value="		u.ID," />
          </node>
          <node concept="l8MVK" id="7tgPrsAl2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAl3" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.UpdatedAt)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAl5" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl6" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAl7" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAl9" role="3cqZAp" />
        <!-- // Delete -->
        <node concept="lc7rE" id="7tgPrsAmn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAma" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmb" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmc" role="lcghm">
            <property role="lacIc" value=") Delete(id int64) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAme" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAmf" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmg" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmh" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmi" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmk" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAml" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAmm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmp" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmr" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAms" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAmt" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAmu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmv" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmx" role="lcghm">
            <property role="lacIc" value="{???-if (schema.hasReferences()) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmz" role="lcghm">
            <property role="lacIc" value="{???-string sn = schema.structName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmB" role="lcghm">
            <property role="lacIc" value="{???-string rn = schema.repoName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmD" role="lcghm">
            <property role="lacIc" value="{???-string tn = schema.name;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAmE" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAmK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmF" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAmG" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmH" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmI" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAmJ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAmL" role="3cqZAp" />
        <!-- // Assign -->
        <node concept="lc7rE" id="7tgPrsAmP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmM" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmN" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmO" role="lcghm">
            <property role="lacIc" value=") Assign(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmR" role="lcghm">
            <property role="lacIc" value="{???-int argIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmT" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmV" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmX" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmZ" role="lcghm">
            <property role="lacIc" value="{???-if (argIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm0" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm3" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm4" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAm5" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm8" role="lcghm">
            <property role="lacIc" value="{???-argIdx = argIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAna" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnc" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <!-- // query -->
        <node concept="lc7rE" id="7tgPrsAnq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnd" role="lcghm">
            <property role="lacIc" value=") (*" />
          </node>
          <node concept="la8eA" id="7tgPrsAne" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnf" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAng" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnh" role="lcghm">
            <property role="lacIc" value="	ur := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAni" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnj" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnl" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnn" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
          <node concept="la8eA" id="7tgPrsAno" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnp" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAns" role="lcghm">
            <property role="lacIc" value="{???-int fkIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnu" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnw" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAny" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnA" role="lcghm">
            <property role="lacIc" value="{???-if (fkIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnB" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnE" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnF" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnI" role="lcghm">
            <property role="lacIc" value="{???-fkIdx = fkIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnK" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnM" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnN" role="lcghm">
            <property role="lacIc" value=") VALUES (" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnQ" role="lcghm">
            <property role="lacIc" value="{???-for (int i = 1; i &lt;= fkIdx; i++) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnS" role="lcghm">
            <property role="lacIc" value="{???-if (i &gt; 1) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnT" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnW" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnX" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
          <node concept="la8eA" id="7tgPrsAnY" role="lcghm">
            <property role="lacIc" value="{???-i}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn1" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn2" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAn4" role="lcghm">
            <property role="lacIc" value="		 ON CONFLICT (" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn7" role="lcghm">
            <property role="lacIc" value="{???-int ckIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn9" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAob" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAod" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAof" role="lcghm">
            <property role="lacIc" value="{???-if (ckIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAog" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoj" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAol" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAok" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAom" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAon" role="lcghm">
            <property role="lacIc" value="{???-ckIdx = ckIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAop" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAor" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAov" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAos" role="lcghm">
            <property role="lacIc" value=") DO NOTHING" />
          </node>
          <node concept="l8MVK" id="7tgPrsAot" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAou" role="lcghm">
            <property role="lacIc" value="		 RETURNING " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAow" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAox" role="lcghm">
            <property role="lacIc" value="{???-int retIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoz" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoB" role="lcghm">
            <property role="lacIc" value="{???-if (retIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoC" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoF" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoH" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoJ" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoK" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoN" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoP" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoR" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoS" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoV" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoX" role="lcghm">
            <property role="lacIc" value="{???-retIdx = retIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoZ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAo2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo0" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
          <node concept="l8MVK" id="7tgPrsAo1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAo3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo4" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAo5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo6" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAo7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo8" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo9" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsApa" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsApb" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsApc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsApe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApf" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAph" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApi" role="lcghm">
            <property role="lacIc" value="	).Scan(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApl" role="lcghm">
            <property role="lacIc" value="{???-int scanIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApn" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApp" role="lcghm">
            <property role="lacIc" value="{???-if (scanIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApq" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAps" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApt" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApv" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApx" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApy" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
          <node concept="la8eA" id="7tgPrsApz" role="lcghm">
            <property role="lacIc" value="{???-fr.pascalName()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApC" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApE" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApG" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApH" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
          <node concept="la8eA" id="7tgPrsApI" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApL" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApN" role="lcghm">
            <property role="lacIc" value="{???-scanIdx = scanIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApP" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApQ" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsApR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApS" role="lcghm">
            <property role="lacIc" value="	return ur, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsApT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApU" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsApV" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsApW" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsApY" role="3cqZAp" />
        <!-- // Remove -->
        <node concept="lc7rE" id="7tgPrsAp2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApZ" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAp0" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAp1" role="lcghm">
            <property role="lacIc" value=") Remove(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAp3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAp4" role="lcghm">
            <property role="lacIc" value="{???-int rmIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAp5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAp6" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAp7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAp8" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAp9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqa" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqc" role="lcghm">
            <property role="lacIc" value="{???-if (rmIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqd" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqg" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqh" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAqi" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAql" role="lcghm">
            <property role="lacIc" value="{???-rmIdx = rmIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqn" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqp" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqq" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqr" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqs" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAqt" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAqu" role="lcghm">
            <property role="lacIc" value=" WHERE " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqx" role="lcghm">
            <property role="lacIc" value="{???-int wIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqz" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqB" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqD" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqF" role="lcghm">
            <property role="lacIc" value="{???-if (wIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqG" role="lcghm">
            <property role="lacIc" value=" AND " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqJ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqL" role="lcghm">
            <property role="lacIc" value="{???-wIdx = wIdx + 1;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqM" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAqN" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
          <node concept="la8eA" id="7tgPrsAqO" role="lcghm">
            <property role="lacIc" value="{???-wIdx}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqR" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqT" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqU" role="lcghm">
            <property role="lacIc" value="`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqX" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqZ" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAq0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq1" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAq4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq2" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAq3" role="lcghm">
            <property role="lacIc" value="{???-fr.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAq5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq6" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAq7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq8" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq9" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAra" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArb" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
          <node concept="l8MVK" id="7tgPrsArc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArd" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAre" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsArf" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsArh" role="3cqZAp" />
        <!-- // Cross-queries -->
        <node concept="lc7rE" id="7tgPrsAri" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArj" role="lcghm">
            <property role="lacIc" value="{???-foreach refA in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArl" role="lcghm">
            <property role="lacIc" value="{???-if (refA.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArn" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAro" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArp" role="lcghm">
            <property role="lacIc" value="{???-foreach refB in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArr" role="lcghm">
            <property role="lacIc" value="{???-if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArt" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAru" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArv" role="lcghm">
            <property role="lacIc" value="{???-string ts = frB.target_schema.structName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArx" role="lcghm">
            <property role="lacIc" value="{???-string tt = frB.target_schema.name;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAry" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsArO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArz" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsArA" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsArB" role="lcghm">
            <property role="lacIc" value=") Get" />
          </node>
          <node concept="la8eA" id="7tgPrsArC" role="lcghm">
            <property role="lacIc" value="{???-ts}" />
          </node>
          <node concept="la8eA" id="7tgPrsArD" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
          <node concept="la8eA" id="7tgPrsArE" role="lcghm">
            <property role="lacIc" value="{???-frA.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsArF" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsArG" role="lcghm">
            <property role="lacIc" value="{???-frA.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsArH" role="lcghm">
            <property role="lacIc" value=" int64) ([]" />
          </node>
          <node concept="la8eA" id="7tgPrsArI" role="lcghm">
            <property role="lacIc" value="{???-ts}" />
          </node>
          <node concept="la8eA" id="7tgPrsArJ" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsArK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArL" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
          <node concept="l8MVK" id="7tgPrsArM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArN" role="lcghm">
            <property role="lacIc" value="		`SELECT t.id" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArQ" role="lcghm">
            <property role="lacIc" value="{???-foreach tf in frB.target_schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArS" role="lcghm">
            <property role="lacIc" value="{???-if (tf.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArU" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) tf;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArV" role="lcghm">
            <property role="lacIc" value=", t." />
          </node>
          <node concept="la8eA" id="7tgPrsArW" role="lcghm">
            <property role="lacIc" value="{???-f.name}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArZ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAr0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAr1" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsE" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAr2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAr3" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAr4" role="lcghm">
            <property role="lacIc" value="{???-tt}" />
          </node>
          <node concept="la8eA" id="7tgPrsAr5" role="lcghm">
            <property role="lacIc" value=" t" />
          </node>
          <node concept="l8MVK" id="7tgPrsAr6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAr7" role="lcghm">
            <property role="lacIc" value="		 INNER JOIN " />
          </node>
          <node concept="la8eA" id="7tgPrsAr8" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAr9" role="lcghm">
            <property role="lacIc" value=" j ON j." />
          </node>
          <node concept="la8eA" id="7tgPrsAsa" role="lcghm">
            <property role="lacIc" value="{???-frB.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAsb" role="lcghm">
            <property role="lacIc" value=" = t.id" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsd" role="lcghm">
            <property role="lacIc" value="		 WHERE j." />
          </node>
          <node concept="la8eA" id="7tgPrsAse" role="lcghm">
            <property role="lacIc" value="{???-frA.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAsf" role="lcghm">
            <property role="lacIc" value=" = $1" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsh" role="lcghm">
            <property role="lacIc" value="		 ORDER BY t.id`, " />
          </node>
          <node concept="la8eA" id="7tgPrsAsi" role="lcghm">
            <property role="lacIc" value="{???-frA.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAsj" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAsk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsl" role="lcghm">
            <property role="lacIc" value="	)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsn" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAso" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsp" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsr" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAss" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAst" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsu" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsv" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
          <node concept="la8eA" id="7tgPrsAsw" role="lcghm">
            <property role="lacIc" value="{???-ts}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsx" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsy" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsA" role="lcghm">
            <property role="lacIc" value="		var item " />
          </node>
          <node concept="la8eA" id="7tgPrsAsB" role="lcghm">
            <property role="lacIc" value="{???-ts}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsD" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;item.ID" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsG" role="lcghm">
            <property role="lacIc" value="{???-foreach tf in frB.target_schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsI" role="lcghm">
            <property role="lacIc" value="{???-if (tf.isInstanceOf(Field)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsK" role="lcghm">
            <property role="lacIc" value="{???-node&lt;Field&gt; f = (Field) tf;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsL" role="lcghm">
            <property role="lacIc" value=", &amp;item." />
          </node>
          <node concept="la8eA" id="7tgPrsAsM" role="lcghm">
            <property role="lacIc" value="{???-f.pascalName()}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsP" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsR" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAs7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsS" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsU" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsW" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsY" role="lcghm">
            <property role="lacIc" value="		items = append(items, item)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAs0" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAs1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAs2" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAs3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAs4" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAs5" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAs6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAs8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAs9" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAta" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtb" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtd" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAte" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtf" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAth" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAti" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtj" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAtk" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — regular schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAtl" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAtm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtn" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAto" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtp" role="lcghm">
            <property role="lacIc" value="{???-if (!(schema.hasReferences())) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtr" role="lcghm">
            <property role="lacIc" value="{???-string sn = schema.structName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAts" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtt" role="lcghm">
            <property role="lacIc" value="{???-string rn = schema.repoName();}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAtu" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAtD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtv" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtw" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtx" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
          <node concept="la8eA" id="7tgPrsAty" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtA" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtB" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAtC" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAtE" role="3cqZAp" />
        <!-- // Create handler -->
        <node concept="lc7rE" id="7tgPrsAuh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtF" role="lcghm">
            <property role="lacIc" value="func handleCreate" />
          </node>
          <node concept="la8eA" id="7tgPrsAtG" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAtH" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAtI" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAtJ" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtL" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtN" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
          <node concept="la8eA" id="7tgPrsAtO" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtQ" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtS" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtU" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtW" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtY" role="lcghm">
            <property role="lacIc" value="		if err := repo.Create(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAt0" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAt1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAt2" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAt3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAt4" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAt5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAt6" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAt7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAt8" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAt9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAua" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAub" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuc" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAud" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAue" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuf" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAug" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAui" role="3cqZAp" />
        <!-- // Get handler -->
        <node concept="lc7rE" id="7tgPrsAuW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuj" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAuk" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAul" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAum" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAun" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuo" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAup" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAur" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAus" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAut" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuu" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuv" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuw" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAux" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuy" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuz" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuA" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuB" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuD" role="lcghm">
            <property role="lacIc" value="		u, err := repo.GetByID(id)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuF" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuH" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusNotFound)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuJ" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuL" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuN" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuO" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuP" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuQ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuR" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuS" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuT" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuU" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAuV" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAuX" role="3cqZAp" />
        <!-- // List handler -->
        <node concept="lc7rE" id="7tgPrsAvp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuY" role="lcghm">
            <property role="lacIc" value="func handleList" />
          </node>
          <node concept="la8eA" id="7tgPrsAuZ" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAu0" role="lcghm">
            <property role="lacIc" value="s(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAu1" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAu2" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAu3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAu4" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAu5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAu6" role="lcghm">
            <property role="lacIc" value="		items, err := repo.List()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAu7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAu8" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAu9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAva" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvc" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAve" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvf" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvg" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvh" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvi" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvk" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvl" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvm" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvn" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAvo" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAvq" role="3cqZAp" />
        <!-- // Update handler -->
        <node concept="lc7rE" id="7tgPrsAwf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvr" role="lcghm">
            <property role="lacIc" value="func handleUpdate" />
          </node>
          <node concept="la8eA" id="7tgPrsAvs" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAvt" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAvu" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAvv" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvw" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvx" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvy" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvz" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvA" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvB" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvD" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvF" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvH" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvJ" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvL" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
          <node concept="la8eA" id="7tgPrsAvM" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvO" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvQ" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvS" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvU" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvW" role="lcghm">
            <property role="lacIc" value="		u.ID = id" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAvY" role="lcghm">
            <property role="lacIc" value="		if err := repo.Update(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAvZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAv0" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAv1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAv2" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAv3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAv4" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAv5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAv6" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAv7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAv8" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAv9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwa" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwc" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwd" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAwe" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAwg" role="3cqZAp" />
        <!-- // Delete handler -->
        <node concept="lc7rE" id="7tgPrsAwQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwh" role="lcghm">
            <property role="lacIc" value="func handleDelete" />
          </node>
          <node concept="la8eA" id="7tgPrsAwi" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAwj" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAwk" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAwl" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwn" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwo" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwp" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwr" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAws" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwt" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwu" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwv" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAww" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwx" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwy" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwz" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwA" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwB" role="lcghm">
            <property role="lacIc" value="		if err := repo.Delete(id); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwD" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwF" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwH" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwJ" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwL" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAwN" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAwO" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAwP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwS" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwU" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAwV" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAwW" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAwX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwY" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAw0" role="lcghm">
            <property role="lacIc" value="{???-if (schema.hasReferences()) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAw1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAw2" role="lcghm">
            <property role="lacIc" value="{???-string sn = schema.structName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAw3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAw4" role="lcghm">
            <property role="lacIc" value="{???-string rn = schema.repoName();}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAw5" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAxf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAw6" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAw7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAw8" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
          <node concept="la8eA" id="7tgPrsAw9" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAxa" role="lcghm">
            <property role="lacIc" value=" (assignments)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxc" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxd" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAxe" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAxg" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAxh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxi" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; firstRef = null;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxk" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; secondRef = null;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxm" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxo" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxq" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxs" role="lcghm">
            <property role="lacIc" value="{???-if (firstRef == null) { firstRef = fr; } }" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxu" role="lcghm">
            <property role="lacIc" value="{???-else if (secondRef == null) { secondRef = fr; } }" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxw" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxy" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAxz" role="3cqZAp" />
        <!-- // Assign handler -->
        <node concept="lc7rE" id="7tgPrsAyv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxA" role="lcghm">
            <property role="lacIc" value="func handleAssign" />
          </node>
          <node concept="la8eA" id="7tgPrsAxB" role="lcghm">
            <property role="lacIc" value="{???-secondRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAxC" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAxD" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAxE" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxG" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxI" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAxJ" role="lcghm">
            <property role="lacIc" value="{???-firstRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAxK" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxM" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxO" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxQ" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxS" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxU" role="lcghm">
            <property role="lacIc" value="		var body Assign" />
          </node>
          <node concept="la8eA" id="7tgPrsAxV" role="lcghm">
            <property role="lacIc" value="{???-secondRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAxW" role="lcghm">
            <property role="lacIc" value="Body" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAxY" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAxZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAx0" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAx1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAx2" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAx3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAx4" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAx5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAx6" role="lcghm">
            <property role="lacIc" value="		ur, err := urRepo.Assign(" />
          </node>
          <node concept="la8eA" id="7tgPrsAx7" role="lcghm">
            <property role="lacIc" value="{???-firstRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAx8" role="lcghm">
            <property role="lacIc" value=", body." />
          </node>
          <node concept="la8eA" id="7tgPrsAx9" role="lcghm">
            <property role="lacIc" value="{???-secondRef.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAya" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyc" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAye" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyf" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyg" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyh" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyi" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyk" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyl" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAym" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyo" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(ur)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyp" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyq" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyr" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAys" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyt" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAyu" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAyw" role="3cqZAp" />
        <!-- // Remove handler -->
        <node concept="lc7rE" id="7tgPrsAzo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyx" role="lcghm">
            <property role="lacIc" value="func handleRemove" />
          </node>
          <node concept="la8eA" id="7tgPrsAyy" role="lcghm">
            <property role="lacIc" value="{???-secondRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAyz" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAyA" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAyB" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyD" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyF" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAyG" role="lcghm">
            <property role="lacIc" value="{???-firstRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAyH" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyJ" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyL" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyN" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyO" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyP" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyQ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyR" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAyS" role="lcghm">
            <property role="lacIc" value="{???-secondRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAyT" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAyU" role="lcghm">
            <property role="lacIc" value="{???-secondRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAyV" role="lcghm">
            <property role="lacIc" value="&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyW" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyX" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAyY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAyZ" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAy0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAy1" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAy2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAy3" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAy4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAy5" role="lcghm">
            <property role="lacIc" value="		if err := urRepo.Remove(" />
          </node>
          <node concept="la8eA" id="7tgPrsAy6" role="lcghm">
            <property role="lacIc" value="{???-firstRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAy7" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAy8" role="lcghm">
            <property role="lacIc" value="{???-secondRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAy9" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAza" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzb" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzd" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAze" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzf" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzh" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzi" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzj" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzl" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzm" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAzn" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAzp" role="3cqZAp" />
        <!-- // Cross-query handlers -->
        <node concept="lc7rE" id="7tgPrsAzq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzr" role="lcghm">
            <property role="lacIc" value="{???-foreach refA in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzt" role="lcghm">
            <property role="lacIc" value="{???-if (refA.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzv" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzx" role="lcghm">
            <property role="lacIc" value="{???-foreach refB in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzz" role="lcghm">
            <property role="lacIc" value="{???-if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzB" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzC" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAzD" role="lcghm">
            <property role="lacIc" value="{???-frA.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAzE" role="lcghm">
            <property role="lacIc" value="{???-frB.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAzF" role="lcghm">
            <property role="lacIc" value="s(urRepo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAzG" role="lcghm">
            <property role="lacIc" value="{???-rn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAzH" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzJ" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzL" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzN" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzO" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzP" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzQ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzR" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzS" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzT" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAzU" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAzV" role="lcghm">
            <property role="lacIc" value="		items, err := urRepo.Get" />
          </node>
          <node concept="la8eA" id="7tgPrsAzW" role="lcghm">
            <property role="lacIc" value="{???-frB.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAzX" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
          <node concept="la8eA" id="7tgPrsAzY" role="lcghm">
            <property role="lacIc" value="{???-frA.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAzZ" role="lcghm">
            <property role="lacIc" value="(id)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAz0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAz1" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAz2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAz3" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAz4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAz5" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAz6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAz7" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAz8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAz9" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAa" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAb" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAd" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAe" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAf" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAg" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAAh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAk" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAm" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAo" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAq" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAs" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAu" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAAv" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Main -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAAw" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsABA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAx" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAy" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAz" role="lcghm">
            <property role="lacIc" value="// Main" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAA" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAB" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAC" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAAD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAE" role="lcghm">
            <property role="lacIc" value="func main() {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAG" role="lcghm">
            <property role="lacIc" value="	dbURL := os.Getenv(&quot;DATABASE_URL&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAI" role="lcghm">
            <property role="lacIc" value="	if dbURL == &quot;&quot; {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAJ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAK" role="lcghm">
            <property role="lacIc" value="		dbURL = &quot;postgres://" />
          </node>
          <node concept="la8eA" id="7tgPrsAAL" role="lcghm">
            <property role="lacIc" value="{???-infra.dbUser}" />
          </node>
          <node concept="la8eA" id="7tgPrsAAM" role="lcghm">
            <property role="lacIc" value=":" />
          </node>
          <node concept="la8eA" id="7tgPrsAAN" role="lcghm">
            <property role="lacIc" value="{???-infra.dbPass}" />
          </node>
          <node concept="la8eA" id="7tgPrsAAO" role="lcghm">
            <property role="lacIc" value="@localhost:5432/" />
          </node>
          <node concept="la8eA" id="7tgPrsAAP" role="lcghm">
            <property role="lacIc" value="{???-infra.dbName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAAQ" role="lcghm">
            <property role="lacIc" value="?sslmode=disable&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAS" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAT" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAAU" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAV" role="lcghm">
            <property role="lacIc" value="	db, err := sql.Open(&quot;postgres&quot;, dbURL)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAW" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAX" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAAY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAAZ" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAA0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAA1" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAA2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAA3" role="lcghm">
            <property role="lacIc" value="	defer db.Close()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAA4" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAA5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAA6" role="lcghm">
            <property role="lacIc" value="	for i := 0; i &lt; 5; i++ {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAA7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAA8" role="lcghm">
            <property role="lacIc" value="		if err = db.Ping(); err == nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAA9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABa" role="lcghm">
            <property role="lacIc" value="			break" />
          </node>
          <node concept="l8MVK" id="7tgPrsABb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABc" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsABd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABe" role="lcghm">
            <property role="lacIc" value="		log.Printf(&quot;DB not ready, retrying... (%d/5)&quot;, i+1)" />
          </node>
          <node concept="l8MVK" id="7tgPrsABf" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABg" role="lcghm">
            <property role="lacIc" value="		time.Sleep(2 * time.Second)" />
          </node>
          <node concept="l8MVK" id="7tgPrsABh" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABi" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsABj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABk" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsABl" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABm" role="lcghm">
            <property role="lacIc" value="		log.Fatal(&quot;DB connection failed:&quot;, err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsABn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABo" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsABp" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsABq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABr" role="lcghm">
            <property role="lacIc" value="	if _, err := db.Exec(migrationSQL); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsABs" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABt" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsABu" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABv" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsABw" role="lcghm" />
          <node concept="la8eA" id="7tgPrsABx" role="lcghm">
            <property role="lacIc" value="	log.Println(&quot;Migration complete&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsABy" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsABz" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsABB" role="3cqZAp" />
        <!-- // Repo instantiation -->
        <node concept="lc7rE" id="7tgPrsABC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABD" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABF" role="lcghm">
            <property role="lacIc" value="{???-if (!(schema.hasReferences())) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABG" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsABH" role="lcghm">
            <property role="lacIc" value="{???-schema.singularName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsABI" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsABJ" role="lcghm">
            <property role="lacIc" value="{???-schema.repoName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsABK" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
          <node concept="l8MVK" id="7tgPrsABL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABO" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABQ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABS" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABU" role="lcghm">
            <property role="lacIc" value="{???-if (schema.hasReferences()) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAB1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABV" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsABW" role="lcghm">
            <property role="lacIc" value="{???-schema.singularName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsABX" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsABY" role="lcghm">
            <property role="lacIc" value="{???-schema.repoName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsABZ" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAB0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAB2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAB3" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAB4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAB5" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAB6" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsACg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAB7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAB8" role="lcghm">
            <property role="lacIc" value="	mux := http.NewServeMux()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAB9" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsACa" role="lcghm" />
          <node concept="la8eA" id="7tgPrsACb" role="lcghm">
            <property role="lacIc" value="	// Swagger UI" />
          </node>
          <node concept="l8MVK" id="7tgPrsACc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsACd" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /swagger/*&quot;, httpSwagger.WrapHandler)" />
          </node>
          <node concept="l8MVK" id="7tgPrsACe" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsACf" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsACh" role="3cqZAp" />
        <!-- // Regular routes -->
        <node concept="lc7rE" id="7tgPrsACi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACj" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACl" role="lcghm">
            <property role="lacIc" value="{???-if (!(schema.hasReferences())) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACn" role="lcghm">
            <property role="lacIc" value="{???-string vr = schema.singularName() + &quot;Repo&quot;;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACp" role="lcghm">
            <property role="lacIc" value="{???-string sn = schema.structName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACr" role="lcghm">
            <property role="lacIc" value="{???-string tn = schema.name;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACs" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
          <node concept="la8eA" id="7tgPrsACt" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACu" role="lcghm">
            <property role="lacIc" value="s" />
          </node>
          <node concept="l8MVK" id="7tgPrsACv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsACw" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
          <node concept="la8eA" id="7tgPrsACx" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACy" role="lcghm">
            <property role="lacIc" value="&quot;, handleCreate" />
          </node>
          <node concept="la8eA" id="7tgPrsACz" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACA" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsACB" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsACC" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsACD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsACE" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsACF" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACG" role="lcghm">
            <property role="lacIc" value="&quot;, handleList" />
          </node>
          <node concept="la8eA" id="7tgPrsACH" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACI" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
          <node concept="la8eA" id="7tgPrsACJ" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsACK" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsACL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsACM" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsACN" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACO" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsACP" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACQ" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsACR" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsACS" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsACT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsACU" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;PUT /" />
          </node>
          <node concept="la8eA" id="7tgPrsACV" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACW" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleUpdate" />
          </node>
          <node concept="la8eA" id="7tgPrsACX" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsACY" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsACZ" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsAC0" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAC1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAC2" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
          <node concept="la8eA" id="7tgPrsAC3" role="lcghm">
            <property role="lacIc" value="{???-tn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAC4" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleDelete" />
          </node>
          <node concept="la8eA" id="7tgPrsAC5" role="lcghm">
            <property role="lacIc" value="{???-sn}" />
          </node>
          <node concept="la8eA" id="7tgPrsAC6" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAC7" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsAC8" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAC9" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsADa" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADd" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADf" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsADg" role="3cqZAp" />
        <!-- // Join routes -->
        <node concept="lc7rE" id="7tgPrsADh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADi" role="lcghm">
            <property role="lacIc" value="{???-foreach schema in model.schemas {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADk" role="lcghm">
            <property role="lacIc" value="{???-if (schema.hasReferences()) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADm" role="lcghm">
            <property role="lacIc" value="{???-string vr = schema.singularName() + &quot;Repo&quot;;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADo" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fRef = null;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADq" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; sRef = null;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADs" role="lcghm">
            <property role="lacIc" value="{???-foreach field in schema.Fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADu" role="lcghm">
            <property role="lacIc" value="{???-if (field.isInstanceOf(FieldRefrence)) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADw" role="lcghm">
            <property role="lacIc" value="{???-node&lt;FieldRefrence&gt; fr = (FieldRefrence) field;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADy" role="lcghm">
            <property role="lacIc" value="{???-if (fRef == null) { fRef = fr; } }" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADA" role="lcghm">
            <property role="lacIc" value="{???-else if (sRef == null) { sRef = fr; } }" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADC" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADE" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADF" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
          <node concept="la8eA" id="7tgPrsADG" role="lcghm">
            <property role="lacIc" value="{???-schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsADH" role="lcghm">
            <property role="lacIc" value=" assignments" />
          </node>
          <node concept="l8MVK" id="7tgPrsADI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsADJ" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
          <node concept="la8eA" id="7tgPrsADK" role="lcghm">
            <property role="lacIc" value="{???-fRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsADL" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsADM" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsADN" role="lcghm">
            <property role="lacIc" value="&quot;, handleAssign" />
          </node>
          <node concept="la8eA" id="7tgPrsADO" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsADP" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsADQ" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsADR" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsADS" role="lcghm" />
          <node concept="la8eA" id="7tgPrsADT" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
          <node concept="la8eA" id="7tgPrsADU" role="lcghm">
            <property role="lacIc" value="{???-fRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsADV" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsADW" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsADX" role="lcghm">
            <property role="lacIc" value="/{" />
          </node>
          <node concept="la8eA" id="7tgPrsADY" role="lcghm">
            <property role="lacIc" value="{???-sRef.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsADZ" role="lcghm">
            <property role="lacIc" value="}&quot;, handleRemove" />
          </node>
          <node concept="la8eA" id="7tgPrsAD0" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAD1" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAD2" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsAD3" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAD4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAD5" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsAD6" role="lcghm">
            <property role="lacIc" value="{???-fRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAD7" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsAD8" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAD9" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAEa" role="lcghm">
            <property role="lacIc" value="{???-fRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEb" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEc" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
          <node concept="la8eA" id="7tgPrsAEd" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEe" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAEf" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAEg" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsAEh" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEi" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsAEj" role="lcghm">
            <property role="lacIc" value="{???-fRef.target_schema.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEk" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAEl" role="lcghm">
            <property role="lacIc" value="{???-sRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEm" role="lcghm">
            <property role="lacIc" value="{???-fRef.target_schema.structName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEn" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
          <node concept="la8eA" id="7tgPrsAEo" role="lcghm">
            <property role="lacIc" value="{???-vr}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEp" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAEq" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAEr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEu" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEw" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAEx" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAEM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEy" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Serving on :" />
          </node>
          <node concept="la8eA" id="7tgPrsAEz" role="lcghm">
            <property role="lacIc" value="{???-infra.port}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEA" role="lcghm">
            <property role="lacIc" value="&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAEB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAEC" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Swagger UI: http://localhost:" />
          </node>
          <node concept="la8eA" id="7tgPrsAED" role="lcghm">
            <property role="lacIc" value="{???-infra.port}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEE" role="lcghm">
            <property role="lacIc" value="/swagger/index.html&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAEF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAEG" role="lcghm">
            <property role="lacIc" value="	log.Fatal(http.ListenAndServe(&quot;:" />
          </node>
          <node concept="la8eA" id="7tgPrsAEH" role="lcghm">
            <property role="lacIc" value="{???-infra.port}" />
          </node>
          <node concept="la8eA" id="7tgPrsAEI" role="lcghm">
            <property role="lacIc" value="&quot;, mux))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAEJ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAEK" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAEL" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAEN" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAEO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEP" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAER" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
      </node>
    </node>
  </node>
</model>
