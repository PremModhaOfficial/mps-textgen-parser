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
  <node concept="WtQ9Q" id="7tgPrsAJV">
    <ref role="WuzLi" to="s1:TODO-CONCEPT-ID" resolve="Code" />
    <node concept="29tfMY" id="7tgPrsAJY" role="29tGrW">
      <node concept="3clFbS" id="7tgPrsAJZ" role="2VODD2">
        <node concept="3cpWs6" id="7tgPrsAJ0" role="3cqZAp">
          <node concept="Xl_RD" id="7tgPrsAJ1" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="7tgPrsAJW" role="11c4hB">
      <node concept="3clFbS" id="7tgPrsAJX" role="2VODD2">
        <!-- VAR: node&lt;Code&gt; n = node; -->
        <!-- VAR: node&lt;Models&gt; model = n.model; -->
        <!-- VAR: node&lt;Infra&gt; infra = n.infra; -->
        <node concept="3clFbH" id="7tgPrsAa5" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAa6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa7" role="lcghm">
            <property role="lacIc" value="package main" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAa8" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAa9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAba" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbd" role="lcghm">
            <property role="lacIc" value="import (" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbe" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbh" role="lcghm">
            <property role="lacIc" value="	&quot;database/sql&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbi" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbl" role="lcghm">
            <property role="lacIc" value="	_ &quot;embed&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbm" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbn" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbp" role="lcghm">
            <property role="lacIc" value="	&quot;encoding/json&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbq" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbt" role="lcghm">
            <property role="lacIc" value="	&quot;fmt&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbu" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbx" role="lcghm">
            <property role="lacIc" value="	&quot;log&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAby" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbB" role="lcghm">
            <property role="lacIc" value="	&quot;net/http&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbC" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbF" role="lcghm">
            <property role="lacIc" value="	&quot;os&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbJ" role="lcghm">
            <property role="lacIc" value="	&quot;strconv&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbN" role="lcghm">
            <property role="lacIc" value="	&quot;time&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbQ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbT" role="lcghm">
            <property role="lacIc" value="	_ &quot;github.com/lib/pq&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbX" role="lcghm">
            <property role="lacIc" value="	httpSwagger &quot;github.com/swaggo/http-swagger&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbY" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAb0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb1" role="lcghm">
            <property role="lacIc" value="	_ &quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAb2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb3" role="lcghm">
            <property role="lacIc" value="▶infra.modulePath◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAb4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb5" role="lcghm">
            <property role="lacIc" value="/docs&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAb6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAb7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAb8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb9" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAca" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcc" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAce" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcf" role="lcghm">
            <property role="lacIc" value="//go:embed user_management_init.sql" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAch" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAci" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcj" role="lcghm">
            <property role="lacIc" value="var migrationSQL string" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAck" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcm" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcn" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAco" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcp" role="lcghm">
            <property role="lacIc" value="// @title         " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcr" role="lcghm">
            <property role="lacIc" value="▶model.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAct" role="lcghm">
            <property role="lacIc" value=" API" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcu" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcx" role="lcghm">
            <property role="lacIc" value="// @version       1.0" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcy" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcB" role="lcghm">
            <property role="lacIc" value="// @description   CRUD service for " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcD" role="lcghm">
            <property role="lacIc" value="▶model.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcE" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcH" role="lcghm">
            <property role="lacIc" value="// @host          localhost:" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcJ" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcN" role="lcghm">
            <property role="lacIc" value="// @BasePath      /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcR" role="lcghm">
            <property role="lacIc" value="// @schemes       http" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcV" role="lcghm">
            <property role="lacIc" value="// @produce       json" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcW" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcZ" role="lcghm">
            <property role="lacIc" value="// @consumes      json" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc0" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAc1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc2" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAc3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc5" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAc7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc9" role="lcghm">
            <property role="lacIc" value="// Models" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAda" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdd" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAde" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdh" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAdi" role="3cqZAp" />
        <!-- // ~~~~ Regular schema structs ~~~~ -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <node concept="3clFbH" id="7tgPrsAdj" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdl" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdn" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdp" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdq" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAds" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdt" role="lcghm">
            <property role="lacIc" value="	ID int64 `json:&quot;id&quot; db:&quot;id&quot; example:&quot;1&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdu" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdv" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAdw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdx" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdz" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdB" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdD" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdF" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdH" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdJ" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdL" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdN" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdP" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAdQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdR" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdV" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAdW" role="3cqZAp" />
        <!-- // MarshalJSON if Sensitive -->
        <!-- VAR: boolean hasSensitive = false; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.Sensitive) { hasSensitive = true; } -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAdX" role="3cqZAp" />
        <!-- // found the sensitive parst -->
        <!-- ─── IF: if (hasSensitive) ─── -->
        <node concept="lc7rE" id="7tgPrsAdY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdZ" role="lcghm">
            <property role="lacIc" value="func (u " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd1" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd3" role="lcghm">
            <property role="lacIc" value=") MarshalJSON() ([]byte, error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd4" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAd5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd7" role="lcghm">
            <property role="lacIc" value="	type Alias " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAd8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd9" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAea" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAec" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAed" role="lcghm">
            <property role="lacIc" value="	return json.Marshal(&amp;struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAee" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAef" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeh" role="lcghm">
            <property role="lacIc" value="		Alias" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAei" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAej" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.Sensitive) ─── -->
        <node concept="lc7rE" id="7tgPrsAek" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAel" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAem" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAen" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAep" role="lcghm">
            <property role="lacIc" value=" string `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAer" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAes" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAet" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeu" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAev" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAew" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAex" role="lcghm">
            <property role="lacIc" value="	}{" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAey" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAez" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeB" role="lcghm">
            <property role="lacIc" value="		Alias: (Alias)(u)," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeC" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeD" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.Sensitive) ─── -->
        <node concept="lc7rE" id="7tgPrsAeE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeF" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeH" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeJ" role="lcghm">
            <property role="lacIc" value=": &quot;[REDACTED]&quot;," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeL" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAeM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeN" role="lcghm">
            <property role="lacIc" value="	})" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeR" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeV" role="lcghm" />
        </node>
        <!-- END IF -->
        <node concept="3clFbH" id="7tgPrsAeW" role="3cqZAp" />
        <node concept="3clFbH" id="7tgPrsAeX" role="3cqZAp" />
        <!-- // Create struct -->
        <node concept="lc7rE" id="7tgPrsAeY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeZ" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe1" role="lcghm">
            <property role="lacIc" value="▶schema.createStructName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe3" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe4" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAe5" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <node concept="lc7rE" id="7tgPrsAe6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe7" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe9" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfb" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfd" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAff" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfh" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfj" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfk" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfl" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAfm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfn" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfo" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfq" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfr" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAfs" role="3cqZAp" />
        <!-- // ~~~~ Join schema structs ~~~~ -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <node concept="lc7rE" id="7tgPrsAft" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfu" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfw" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfy" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfz" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfA" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAfB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfC" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfE" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfG" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfI" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfK" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfM" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfO" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfQ" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAfR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfS" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfU" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfW" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfY" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf0" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf2" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf4" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf6" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf8" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf9" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAga" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAgb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgc" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAge" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgf" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgg" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgh" role="3cqZAp" />
        <!-- VAR: int refCount = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- ─── IF: if (refCount == 1) ─── -->
        <node concept="lc7rE" id="7tgPrsAgi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgj" role="lcghm">
            <property role="lacIc" value="type Assign" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgl" role="lcghm">
            <property role="lacIc" value="▶fr.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgn" role="lcghm">
            <property role="lacIc" value="Body struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgo" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgr" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgt" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgv" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgx" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgz" role="lcghm">
            <property role="lacIc" value="&quot; binding:&quot;required&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgA" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgD" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgE" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgH" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- ASSIGN: refCount = refCount + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAgI" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — regular -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAgJ" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAgK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgL" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgM" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgP" role="lcghm">
            <property role="lacIc" value="// Repositories" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgQ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgT" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgW" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgX" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgY" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="3clFbH" id="7tgPrsAgZ" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAg0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg1" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg3" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg5" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAg7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAg8" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAg9" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAha" role="3cqZAp" />
        <!-- // Create -->
        <node concept="lc7rE" id="7tgPrsAhb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhc" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhe" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhg" role="lcghm">
            <property role="lacIc" value=") Create(u *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhi" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhk" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAhm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAho" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhp" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAhq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhs" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAht" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhu" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhw" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <!-- VAR: int idx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <!-- IF: if (idx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAhx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhy" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAhz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhA" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- ASSIGN: idx = idx + 1; -->
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAhB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhC" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhD" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAhE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhG" role="lcghm">
            <property role="lacIc" value="		 VALUES (" />
          </node>
        </node>
        <!-- ═══ FOR: for (int i = 1; i &lt;= idx; i++) ═══ -->
        <!-- IF: if (i &gt; 1) -->
        <node concept="lc7rE" id="7tgPrsAhH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhI" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAhJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhK" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhM" role="lcghm">
            <property role="lacIc" value="▶i◀" />
          </node>
        </node>
        <!-- END FOR -->
        <node concept="lc7rE" id="7tgPrsAhN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhO" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAhQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhS" role="lcghm">
            <property role="lacIc" value="		 RETURNING id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType == DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAhT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhU" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhW" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAhX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhY" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAh0" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAh1" role="3cqZAp" />
        <!-- // non timestapm -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType != DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAh2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh3" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAh4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh5" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAh6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh7" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAh8" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAh9" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAia" role="3cqZAp" />
        <!-- // scanning -->
        <node concept="lc7rE" id="7tgPrsAib" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAic" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType == DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAid" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAie" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAif" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAig" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAih" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAii" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAij" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAik" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAil" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAim" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAin" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAio" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAip" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAiq" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAir" role="3cqZAp" />
        <!-- // GetByID -->
        <node concept="lc7rE" id="7tgPrsAis" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAit" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiv" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAix" role="lcghm">
            <property role="lacIc" value=") GetByID(id int64) (*" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiz" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiB" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiC" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAiD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiF" role="lcghm">
            <property role="lacIc" value="	u := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiH" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiJ" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAiL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiN" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAiP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiR" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAiS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiT" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiV" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAiW" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAiX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiZ" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi1" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi3" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi4" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAi5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAi6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi7" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAi8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi9" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAja" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjb" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAjc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjd" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAje" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjh" role="lcghm">
            <property role="lacIc" value="	return u, err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAji" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjl" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjm" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjn" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjo" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjp" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAjq" role="3cqZAp" />
        <!-- // List -->
        <node concept="lc7rE" id="7tgPrsAjr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjs" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAju" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjw" role="lcghm">
            <property role="lacIc" value=") List() ([]" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjy" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjA" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjE" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjF" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjI" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAjJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjK" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjM" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAjN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjQ" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjS" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjU" role="lcghm">
            <property role="lacIc" value=" ORDER BY id`)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjV" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjY" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAj0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAj1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj2" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAj3" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAj4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAj5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj6" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAj7" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAj8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAj9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAka" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkb" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAke" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkg" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkh" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAki" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkk" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAko" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkq" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkr" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAks" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAku" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAkv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkw" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAky" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAkz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkA" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkE" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkF" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkI" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkJ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkM" role="lcghm">
            <property role="lacIc" value="		items = append(items, u)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkQ" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkR" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkS" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkU" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkV" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkY" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAk0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAk1" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAk2" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAk3" role="3cqZAp" />
        <!-- // Update -->
        <node concept="lc7rE" id="7tgPrsAk4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk5" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAk6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk7" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAk8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk9" role="lcghm">
            <property role="lacIc" value=") Update(u *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAla" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlb" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAld" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAle" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlh" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAli" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAll" role="lcghm">
            <property role="lacIc" value="		`UPDATE " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAln" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlp" role="lcghm">
            <property role="lacIc" value=" SET " />
          </node>
        </node>
        <!-- VAR: int uidx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <!-- IF: if (uidx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAlq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlr" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: uidx = uidx + 1; -->
        <node concept="lc7rE" id="7tgPrsAls" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlt" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlv" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlx" role="lcghm">
            <property role="lacIc" value="▶uidx◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- VAR: int nextParam = uidx + 1; -->
        <node concept="lc7rE" id="7tgPrsAly" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlz" role="lcghm">
            <property role="lacIc" value=", updated_at = NOW()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlA" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlD" role="lcghm">
            <property role="lacIc" value="		 WHERE id = $" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlF" role="lcghm">
            <property role="lacIc" value="▶nextParam◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlJ" role="lcghm">
            <property role="lacIc" value="		 RETURNING updated_at`," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlL" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType != DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAlM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlN" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlP" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlR" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlT" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAlU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlV" role="lcghm">
            <property role="lacIc" value="		u.ID," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlW" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlZ" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.UpdatedAt)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAl0" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAl1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAl2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAl3" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAl4" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAl5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAl6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAl7" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAl8" role="3cqZAp" />
        <!-- // Delete -->
        <node concept="lc7rE" id="7tgPrsAl9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAma" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmc" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAme" role="lcghm">
            <property role="lacIc" value=") Delete(id int64) error {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmf" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmi" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmk" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAml" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmm" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmn" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmq" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmr" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAms" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmu" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmv" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmx" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmy" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAmz" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAmA" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="3clFbH" id="7tgPrsAmB" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAmC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmD" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmF" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmH" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmL" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAmM" role="3cqZAp" />
        <!-- // Assign -->
        <node concept="lc7rE" id="7tgPrsAmN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmO" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmQ" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmS" role="lcghm">
            <property role="lacIc" value=") Assign(" />
          </node>
        </node>
        <!-- VAR: int argIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (argIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAmT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmU" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAmV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmW" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmY" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <!-- ASSIGN: argIdx = argIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- // query -->
        <node concept="lc7rE" id="7tgPrsAmZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm0" role="lcghm">
            <property role="lacIc" value=") (*" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm2" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm4" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm5" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAm6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAm7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm8" role="lcghm">
            <property role="lacIc" value="	ur := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAm9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAna" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnc" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAne" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAng" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnh" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAni" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnk" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnm" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAno" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <!-- VAR: int fkIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (fkIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAnp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnq" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAnr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAns" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- ASSIGN: fkIdx = fkIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAnt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnu" role="lcghm">
            <property role="lacIc" value=") VALUES (" />
          </node>
        </node>
        <!-- ═══ FOR: for (int i = 1; i &lt;= fkIdx; i++) ═══ -->
        <!-- IF: if (i &gt; 1) -->
        <node concept="lc7rE" id="7tgPrsAnv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnw" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAnx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAny" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnA" role="lcghm">
            <property role="lacIc" value="▶i◀" />
          </node>
        </node>
        <!-- END FOR -->
        <node concept="lc7rE" id="7tgPrsAnB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnC" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnD" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAnE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnG" role="lcghm">
            <property role="lacIc" value="		 ON CONFLICT (" />
          </node>
        </node>
        <!-- VAR: int ckIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (ckIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAnH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnI" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAnJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnK" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- ASSIGN: ckIdx = ckIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAnL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnM" role="lcghm">
            <property role="lacIc" value=") DO NOTHING" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAnO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnQ" role="lcghm">
            <property role="lacIc" value="		 RETURNING " />
          </node>
        </node>
        <!-- VAR: int retIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- IF: if (retIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAnR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnS" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAnT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnU" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAnV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnW" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: retIdx = retIdx + 1; -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAnX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnY" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAn0" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAn1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn2" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn4" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn6" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn7" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAn8" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAn9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoa" role="lcghm">
            <property role="lacIc" value="	).Scan(" />
          </node>
        </node>
        <!-- VAR: int scanIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- IF: if (scanIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAob" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoc" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAod" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoe" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAof" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAog" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAoh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoi" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAok" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: scanIdx = scanIdx + 1; -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAol" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAom" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAon" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAoo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAop" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoq" role="lcghm">
            <property role="lacIc" value="	return ur, err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAor" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAos" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAot" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAou" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAov" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAow" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAox" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAoy" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAoz" role="3cqZAp" />
        <!-- // Remove -->
        <node concept="lc7rE" id="7tgPrsAoA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoB" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoD" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoF" role="lcghm">
            <property role="lacIc" value=") Remove(" />
          </node>
        </node>
        <!-- VAR: int rmIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (rmIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAoG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoH" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAoI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoJ" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoL" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <!-- ASSIGN: rmIdx = rmIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAoM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoN" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAoP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoR" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoT" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoV" role="lcghm">
            <property role="lacIc" value=" WHERE " />
          </node>
        </node>
        <!-- VAR: int wIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (wIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAoW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoX" role="lcghm">
            <property role="lacIc" value=" AND " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: wIdx = wIdx + 1; -->
        <node concept="lc7rE" id="7tgPrsAoY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoZ" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAo0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo1" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAo2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo3" role="lcghm">
            <property role="lacIc" value="▶wIdx◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAo4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo5" role="lcghm">
            <property role="lacIc" value="`" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAo6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo7" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAo8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo9" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsApa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApb" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApc" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsApd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsApe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApf" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAph" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsApi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApj" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApk" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsApl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsApm" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsApn" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsApo" role="3cqZAp" />
        <!-- // Cross-queries -->
        <!-- ═══ FOREACH: foreach refA in schema.Fields ═══ -->
        <!-- ─── IF: if (refA.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA; -->
        <!-- ═══ FOREACH: foreach refB in schema.Fields ═══ -->
        <!-- ─── IF: if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB; -->
        <!-- VAR: string ts = frB.target_schema.structName(); -->
        <!-- VAR: string tt = frB.target_schema.name; -->
        <node concept="3clFbH" id="7tgPrsApp" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsApq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApr" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAps" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApt" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApv" role="lcghm">
            <property role="lacIc" value=") Get" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApx" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApz" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApB" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApD" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApF" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApH" role="lcghm">
            <property role="lacIc" value=" int64) ([]" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApJ" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApL" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApM" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsApN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsApO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApP" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApQ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsApR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsApS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApT" role="lcghm">
            <property role="lacIc" value="		`SELECT t.id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach tf in frB.target_schema.Fields ═══ -->
        <!-- ─── IF: if (tf.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) tf; -->
        <node concept="lc7rE" id="7tgPrsApU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApV" role="lcghm">
            <property role="lacIc" value=", t." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApX" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsApY" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsApZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAp0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAp1" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAp2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAp3" role="lcghm">
            <property role="lacIc" value="▶tt◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAp4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAp5" role="lcghm">
            <property role="lacIc" value=" t" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAp6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAp7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAp8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAp9" role="lcghm">
            <property role="lacIc" value="		 INNER JOIN " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqb" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqd" role="lcghm">
            <property role="lacIc" value=" j ON j." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqf" role="lcghm">
            <property role="lacIc" value="▶frB.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqh" role="lcghm">
            <property role="lacIc" value=" = t.id" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqi" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAql" role="lcghm">
            <property role="lacIc" value="		 WHERE j." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqn" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqp" role="lcghm">
            <property role="lacIc" value=" = $1" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqq" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqt" role="lcghm">
            <property role="lacIc" value="		 ORDER BY t.id`, " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqv" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqx" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqy" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqB" role="lcghm">
            <property role="lacIc" value="	)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqC" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqF" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqJ" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqN" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqR" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAqU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqV" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqX" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAqY" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAqZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAq0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq1" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAq2" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAq3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAq4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq5" role="lcghm">
            <property role="lacIc" value="		var item " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAq6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAq7" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAq8" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAq9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAra" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArb" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;item.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach tf in frB.target_schema.Fields ═══ -->
        <!-- ─── IF: if (tf.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) tf; -->
        <node concept="lc7rE" id="7tgPrsArc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArd" role="lcghm">
            <property role="lacIc" value=", &amp;item." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAre" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArf" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsArg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArh" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAri" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArl" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArm" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArn" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAro" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArp" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArq" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArt" role="lcghm">
            <property role="lacIc" value="		items = append(items, item)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAru" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArx" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAry" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArB" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArC" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArF" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArJ" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsArK" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — regular schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsArL" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <node concept="3clFbH" id="7tgPrsArM" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsArN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArO" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArS" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArU" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArV" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsArW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsArX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArY" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsArZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAr0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAr1" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAr2" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAr3" role="3cqZAp" />
        <!-- // Create handler -->
        <node concept="lc7rE" id="7tgPrsAr4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAr5" role="lcghm">
            <property role="lacIc" value="func handleCreate" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAr6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAr7" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAr8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAr9" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsb" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsd" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAse" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsh" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsi" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsl" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsn" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAso" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsr" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAss" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAst" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsv" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsw" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsx" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsz" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsA" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsD" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsE" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsH" role="lcghm">
            <property role="lacIc" value="		if err := repo.Create(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsL" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsM" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsP" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsQ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsT" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAsW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsX" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAsY" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAs0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAs1" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAs2" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAs3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAs4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAs5" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAs6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAs7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAs8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAs9" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAta" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtd" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAte" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAth" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAti" role="3cqZAp" />
        <!-- // Get handler -->
        <node concept="lc7rE" id="7tgPrsAtj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtk" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtm" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAto" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtq" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAts" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtt" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtw" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtx" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAty" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtA" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtE" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtF" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtI" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtJ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtM" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtQ" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtR" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtS" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtU" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtV" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAtW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAtX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtY" role="lcghm">
            <property role="lacIc" value="		u, err := repo.GetByID(id)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAtZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAt0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAt1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAt2" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAt3" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAt4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAt5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAt6" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusNotFound)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAt7" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAt8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAt9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAua" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAub" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAud" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAue" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuf" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAug" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAuh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAui" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAul" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAum" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAun" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAup" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuq" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAur" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAus" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAut" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuu" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuv" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAux" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuy" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAuz" role="3cqZAp" />
        <!-- // List handler -->
        <node concept="lc7rE" id="7tgPrsAuA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuB" role="lcghm">
            <property role="lacIc" value="func handleList" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuD" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuF" role="lcghm">
            <property role="lacIc" value="s(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuH" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuJ" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAuM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuN" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAuQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuR" role="lcghm">
            <property role="lacIc" value="		items, err := repo.List()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAuU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuV" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAuW" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAuX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAuY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuZ" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAu0" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAu1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAu2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAu3" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAu4" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAu5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAu6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAu7" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAu8" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAu9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAva" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvb" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvc" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAve" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvf" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvj" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvk" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvn" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvo" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvq" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvr" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAvs" role="3cqZAp" />
        <!-- // Update handler -->
        <node concept="lc7rE" id="7tgPrsAvt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvu" role="lcghm">
            <property role="lacIc" value="func handleUpdate" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvw" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvy" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvA" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvC" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvD" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvG" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvH" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvI" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvK" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvL" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvO" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvS" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvT" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAvW" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAvX" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAvY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAvZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAv0" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAv1" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAv2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAv3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAv4" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAv5" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAv6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAv7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAv8" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAv9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwa" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwb" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwe" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwf" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwi" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwm" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwn" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwq" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwr" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAws" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwu" role="lcghm">
            <property role="lacIc" value="		u.ID = id" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwv" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAww" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwy" role="lcghm">
            <property role="lacIc" value="		if err := repo.Update(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwz" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwA" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwC" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwD" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwG" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwH" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwI" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwK" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwL" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwO" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwS" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwT" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAwW" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAwX" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAwY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAwZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAw0" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAw1" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAw2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAw3" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAw4" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAw5" role="3cqZAp" />
        <!-- // Delete handler -->
        <node concept="lc7rE" id="7tgPrsAw6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAw7" role="lcghm">
            <property role="lacIc" value="func handleDelete" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAw8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAw9" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxb" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxd" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxf" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxj" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxk" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxn" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxo" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxr" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxs" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxt" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxv" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxw" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxx" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxz" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxA" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxD" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxE" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxH" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxL" role="lcghm">
            <property role="lacIc" value="		if err := repo.Delete(id); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxM" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxP" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxQ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxT" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAxW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAxX" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAxY" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAxZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAx0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAx1" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAx2" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAx3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAx4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAx5" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAx6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAx7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAx8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAx9" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAya" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAyc" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyd" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAye" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAyf" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <node concept="3clFbH" id="7tgPrsAyg" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAyh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyi" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAyl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAym" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyo" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyq" role="lcghm">
            <property role="lacIc" value=" (assignments)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyr" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAys" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAyt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyu" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyv" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAyx" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyy" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAyz" role="3cqZAp" />
        <!-- VAR: node&lt;FieldRefrence&gt; firstRef = null; -->
        <!-- VAR: node&lt;FieldRefrence&gt; secondRef = null; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (firstRef == null) { firstRef = fr; } -->
        <!-- IF: else if (secondRef == null) { secondRef = fr; } -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAyA" role="3cqZAp" />
        <!-- // Assign handler -->
        <node concept="lc7rE" id="7tgPrsAyB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyC" role="lcghm">
            <property role="lacIc" value="func handleAssign" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyE" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyG" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyI" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyK" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyL" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAyN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyO" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAyR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyS" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyU" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAyW" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAyX" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAyY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAyZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAy0" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAy1" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAy2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAy3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAy4" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAy5" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAy6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAy7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAy8" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAy9" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAza" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzc" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAze" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzg" role="lcghm">
            <property role="lacIc" value="		var body Assign" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzi" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzk" role="lcghm">
            <property role="lacIc" value="Body" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzo" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzp" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzs" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzt" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzw" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzx" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzy" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzA" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzE" role="lcghm">
            <property role="lacIc" value="		ur, err := urRepo.Assign(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzG" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzI" role="lcghm">
            <property role="lacIc" value=", body." />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzK" role="lcghm">
            <property role="lacIc" value="▶secondRef.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzM" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzQ" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzR" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzS" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzU" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzV" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAzW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAzX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAzY" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAzZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAz0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAz1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAz2" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAz3" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAz4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAz5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAz6" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAz7" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAz8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAz9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAa" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAb" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAe" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(ur)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAf" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAi" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAm" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAn" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAp" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAq" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAAr" role="3cqZAp" />
        <!-- // Remove handler -->
        <node concept="lc7rE" id="7tgPrsAAs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAt" role="lcghm">
            <property role="lacIc" value="func handleRemove" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAv" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAx" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAz" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAB" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAC" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAF" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAJ" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAL" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAN" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAR" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAV" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAAW" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAAX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAAY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAAZ" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAA0" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAA1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAA2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAA3" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAA4" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAA5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAA6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAA7" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAA8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAA9" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABb" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABd" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABf" role="lcghm">
            <property role="lacIc" value="&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABj" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABk" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABn" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABo" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABr" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABs" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABt" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABv" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABw" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABx" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABz" role="lcghm">
            <property role="lacIc" value="		if err := urRepo.Remove(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABB" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABD" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABF" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABH" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABL" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABM" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABP" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABQ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABT" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsABW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsABX" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsABY" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsABZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAB0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAB1" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAB2" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAB3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAB4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAB5" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAB6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAB7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAB8" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAB9" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsACa" role="3cqZAp" />
        <!-- // Cross-query handlers -->
        <!-- ═══ FOREACH: foreach refA in schema.Fields ═══ -->
        <!-- ─── IF: if (refA.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA; -->
        <!-- ═══ FOREACH: foreach refB in schema.Fields ═══ -->
        <!-- ─── IF: if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB; -->
        <node concept="lc7rE" id="7tgPrsACb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACc" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACe" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACg" role="lcghm">
            <property role="lacIc" value="▶frB.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACi" role="lcghm">
            <property role="lacIc" value="s(urRepo *" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACk" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACm" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACn" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACq" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACr" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACs" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACu" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACv" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACy" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACz" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACA" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACC" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACD" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACG" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACH" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACI" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACK" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACL" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACO" role="lcghm">
            <property role="lacIc" value="		items, err := urRepo.Get" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACQ" role="lcghm">
            <property role="lacIc" value="▶frB.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACS" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACU" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsACW" role="lcghm">
            <property role="lacIc" value="(id)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsACX" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsACY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsACZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAC0" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAC1" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAC2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAC3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAC4" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAC5" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAC6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAC7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAC8" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAC9" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADa" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADc" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADe" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADg" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADh" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADi" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADk" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADo" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADp" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADs" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADt" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADv" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADw" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsADx" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Main -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsADy" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsADz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADA" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADE" role="lcghm">
            <property role="lacIc" value="// Main" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADF" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADI" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADJ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADL" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADO" role="lcghm">
            <property role="lacIc" value="func main() {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADS" role="lcghm">
            <property role="lacIc" value="	dbURL := os.Getenv(&quot;DATABASE_URL&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADT" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsADW" role="lcghm">
            <property role="lacIc" value="	if dbURL == &quot;&quot; {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsADX" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsADY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsADZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAD0" role="lcghm">
            <property role="lacIc" value="		dbURL = &quot;postgres://" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAD1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAD2" role="lcghm">
            <property role="lacIc" value="▶infra.dbUser◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAD3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAD4" role="lcghm">
            <property role="lacIc" value=":" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAD5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAD6" role="lcghm">
            <property role="lacIc" value="▶infra.dbPass◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAD7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAD8" role="lcghm">
            <property role="lacIc" value="@localhost:5432/" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAD9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEa" role="lcghm">
            <property role="lacIc" value="▶infra.dbName◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEc" role="lcghm">
            <property role="lacIc" value="?sslmode=disable&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEe" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEg" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEh" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEi" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEm" role="lcghm">
            <property role="lacIc" value="	db, err := sql.Open(&quot;postgres&quot;, dbURL)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEn" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEq" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEr" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEs" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEu" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEv" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEy" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEz" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEA" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEC" role="lcghm">
            <property role="lacIc" value="	defer db.Close()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAED" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEF" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEI" role="lcghm">
            <property role="lacIc" value="	for i := 0; i &lt; 5; i++ {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEJ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEM" role="lcghm">
            <property role="lacIc" value="		if err = db.Ping(); err == nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEQ" role="lcghm">
            <property role="lacIc" value="			break" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAER" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAES" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAET" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEU" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEV" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAEW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAEX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAEY" role="lcghm">
            <property role="lacIc" value="		log.Printf(&quot;DB not ready, retrying... (%d/5)&quot;, i+1)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAEZ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAE0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAE1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAE2" role="lcghm">
            <property role="lacIc" value="		time.Sleep(2 * time.Second)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAE3" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAE4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAE5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAE6" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAE7" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAE8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAE9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFa" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFb" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFe" role="lcghm">
            <property role="lacIc" value="		log.Fatal(&quot;DB connection failed:&quot;, err)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFf" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFi" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFo" role="lcghm">
            <property role="lacIc" value="	if _, err := db.Exec(migrationSQL); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFp" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFs" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFt" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFw" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFx" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFy" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFA" role="lcghm">
            <property role="lacIc" value="	log.Println(&quot;Migration complete&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAFD" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFE" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAFF" role="3cqZAp" />
        <!-- // Repo instantiation -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <node concept="lc7rE" id="7tgPrsAFG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFH" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFJ" role="lcghm">
            <property role="lacIc" value="▶schema.singularName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFL" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFN" role="lcghm">
            <property role="lacIc" value="▶schema.repoName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFP" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFQ" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAFR" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <node concept="lc7rE" id="7tgPrsAFS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFT" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFV" role="lcghm">
            <property role="lacIc" value="▶schema.singularName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFX" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAFY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAFZ" role="lcghm">
            <property role="lacIc" value="▶schema.repoName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAF0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAF1" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAF2" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAF3" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAF4" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAF5" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAF6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAF7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAF8" role="lcghm">
            <property role="lacIc" value="	mux := http.NewServeMux()" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAF9" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAGa" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAGb" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAGc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAGd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGe" role="lcghm">
            <property role="lacIc" value="	// Swagger UI" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGf" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAGg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAGh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGi" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /swagger/*&quot;, httpSwagger.WrapHandler)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGj" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAGk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAGl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAGm" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAGn" role="3cqZAp" />
        <!-- // Regular routes -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string vr = schema.singularName() + "Repo"; -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="lc7rE" id="7tgPrsAGo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGp" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGr" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGt" role="lcghm">
            <property role="lacIc" value="s" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGu" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAGv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAGw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGx" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGz" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGB" role="lcghm">
            <property role="lacIc" value="&quot;, handleCreate" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGD" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGF" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGH" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGJ" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAGL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAGM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGN" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGP" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGR" role="lcghm">
            <property role="lacIc" value="&quot;, handleList" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGT" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGV" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGX" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAGY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAGZ" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAG0" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAG1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAG2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAG3" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAG4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAG5" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAG6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAG7" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAG8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAG9" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHb" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHd" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHf" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAHh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAHi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHj" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;PUT /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHl" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHn" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleUpdate" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHp" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHr" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHt" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHv" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHw" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAHx" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAHy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHz" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHB" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHD" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleDelete" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHF" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHH" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHJ" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHL" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHM" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAHN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAHO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAHP" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAHQ" role="3cqZAp" />
        <!-- // Join routes -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <!-- VAR: string vr = schema.singularName() + "Repo"; -->
        <!-- VAR: node&lt;FieldRefrence&gt; fRef = null; -->
        <!-- VAR: node&lt;FieldRefrence&gt; sRef = null; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (fRef == null) { fRef = fr; } -->
        <!-- IF: else if (sRef == null) { sRef = fr; } -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAHR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHS" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHU" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAHW" role="lcghm">
            <property role="lacIc" value=" assignments" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAHX" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAHY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAHZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAH0" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAH1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAH2" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAH3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAH4" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAH5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAH6" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAH7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAH8" role="lcghm">
            <property role="lacIc" value="&quot;, handleAssign" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAH9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIa" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIc" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAId" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIe" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIg" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIh" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAIi" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAIj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIk" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIm" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIo" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIq" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIs" role="lcghm">
            <property role="lacIc" value="/{" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIu" role="lcghm">
            <property role="lacIc" value="▶sRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIw" role="lcghm">
            <property role="lacIc" value="}&quot;, handleRemove" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIy" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIA" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIC" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAID" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIE" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIF" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAIG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAIH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAII" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIK" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIM" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIO" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIQ" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIS" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIU" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIW" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAIY" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAIZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAI0" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAI1" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAI2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAI3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAI4" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAI5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAI6" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAI7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAI8" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAI9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJa" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJc" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJe" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJg" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJi" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJk" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJm" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJn" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAJo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAJp" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAJq" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAJr" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAJs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJt" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Serving on :" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJv" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJx" role="lcghm">
            <property role="lacIc" value="&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJy" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAJz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAJA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJB" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Swagger UI: http://localhost:" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJD" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJF" role="lcghm">
            <property role="lacIc" value="/swagger/index.html&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAJH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAJI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJJ" role="lcghm">
            <property role="lacIc" value="	log.Fatal(http.ListenAndServe(&quot;:" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJL" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJN" role="lcghm">
            <property role="lacIc" value="&quot;, mux))" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAJP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAJQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAJR" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAJS" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAJT" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAJU" role="3cqZAp" />
        <!-- END ? -->
        <!-- END ? -->
      </node>
    </node>
  </node>
</model>
