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
  <node concept="WtQ9Q" id="7tgPrsAuA">
    <ref role="WuzLi" to="s1:TODO-CONCEPT-ID" resolve="Code" />
    <node concept="29tfMY" id="7tgPrsAuD" role="29tGrW">
      <node concept="3clFbS" id="7tgPrsAuE" role="2VODD2">
        <node concept="3cpWs6" id="7tgPrsAuF" role="3cqZAp">
          <node concept="Xl_RD" id="7tgPrsAuG" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="7tgPrsAuB" role="11c4hB">
      <node concept="3clFbS" id="7tgPrsAuC" role="2VODD2">
        <!-- VAR: node&lt;Code&gt; n = node; -->
        <!-- VAR: node&lt;Models&gt; model = n.model; -->
        <!-- VAR: node&lt;Infra&gt; infra = n.infra; -->
        <node concept="3clFbH" id="7tgPrsAa5" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa6" role="lcghm">
            <property role="lacIc" value="package main" />
          </node>
          <node concept="l8MVK" id="7tgPrsAa7" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAa8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAa9" role="lcghm">
            <property role="lacIc" value="import (" />
          </node>
          <node concept="l8MVK" id="7tgPrsAba" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbb" role="lcghm">
            <property role="lacIc" value="	&quot;database/sql&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbd" role="lcghm">
            <property role="lacIc" value="	_ &quot;embed&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbe" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbf" role="lcghm">
            <property role="lacIc" value="	&quot;encoding/json&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbh" role="lcghm">
            <property role="lacIc" value="	&quot;fmt&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbi" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbj" role="lcghm">
            <property role="lacIc" value="	&quot;log&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbl" role="lcghm">
            <property role="lacIc" value="	&quot;net/http&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbn" role="lcghm">
            <property role="lacIc" value="	&quot;os&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbo" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbp" role="lcghm">
            <property role="lacIc" value="	&quot;strconv&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbr" role="lcghm">
            <property role="lacIc" value="	&quot;time&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbs" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbt" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbu" role="lcghm">
            <property role="lacIc" value="	_ &quot;github.com/lib/pq&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbw" role="lcghm">
            <property role="lacIc" value="	httpSwagger &quot;github.com/swaggo/http-swagger&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbx" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAby" role="lcghm">
            <property role="lacIc" value="	_ &quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAbz" role="lcghm">
            <property role="lacIc" value="▶infra.modulePath◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAbA" role="lcghm">
            <property role="lacIc" value="/docs&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbC" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbD" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbF" role="lcghm">
            <property role="lacIc" value="//go:embed user_management_init.sql" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbH" role="lcghm">
            <property role="lacIc" value="var migrationSQL string" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbI" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbJ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbK" role="lcghm">
            <property role="lacIc" value="// @title         " />
          </node>
          <node concept="la8eA" id="7tgPrsAbL" role="lcghm">
            <property role="lacIc" value="▶model.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAbM" role="lcghm">
            <property role="lacIc" value=" API" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbO" role="lcghm">
            <property role="lacIc" value="// @version       1.0" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbQ" role="lcghm">
            <property role="lacIc" value="// @description   CRUD service for " />
          </node>
          <node concept="la8eA" id="7tgPrsAbR" role="lcghm">
            <property role="lacIc" value="▶model.name◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbS" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbT" role="lcghm">
            <property role="lacIc" value="// @host          localhost:" />
          </node>
          <node concept="la8eA" id="7tgPrsAbU" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbW" role="lcghm">
            <property role="lacIc" value="// @BasePath      /" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAbY" role="lcghm">
            <property role="lacIc" value="// @schemes       http" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb0" role="lcghm">
            <property role="lacIc" value="// @produce       json" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb2" role="lcghm">
            <property role="lacIc" value="// @consumes      json" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb3" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAb4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb5" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb7" role="lcghm">
            <property role="lacIc" value="// Models" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAb9" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAca" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAcb" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAcd" role="3cqZAp" />
        <!-- // ~~~~ Regular schema structs ~~~~ -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <node concept="3clFbH" id="7tgPrsAce" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcf" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAcg" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAch" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAci" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcj" role="lcghm">
            <property role="lacIc" value="	ID int64 `json:&quot;id&quot; db:&quot;id&quot; example:&quot;1&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAck" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAcw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcm" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcn" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAco" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAcp" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAcq" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcr" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAcs" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAct" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAcu" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcv" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAcA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcx" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcy" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAcz" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAcB" role="3cqZAp" />
        <!-- // MarshalJSON if Sensitive -->
        <!-- VAR: boolean hasSensitive = false; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.Sensitive) { hasSensitive = true; } -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAcC" role="3cqZAp" />
        <!-- // found the sensitive parst -->
        <!-- ─── IF: if (hasSensitive) ─── -->
        <node concept="lc7rE" id="7tgPrsAcO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcD" role="lcghm">
            <property role="lacIc" value="func (u " />
          </node>
          <node concept="la8eA" id="7tgPrsAcE" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAcF" role="lcghm">
            <property role="lacIc" value=") MarshalJSON() ([]byte, error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcH" role="lcghm">
            <property role="lacIc" value="	type Alias " />
          </node>
          <node concept="la8eA" id="7tgPrsAcI" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcJ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcK" role="lcghm">
            <property role="lacIc" value="	return json.Marshal(&amp;struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcM" role="lcghm">
            <property role="lacIc" value="		Alias" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcN" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.Sensitive) ─── -->
        <node concept="lc7rE" id="7tgPrsAcV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcP" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAcQ" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAcR" role="lcghm">
            <property role="lacIc" value=" string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcS" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAcT" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcU" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAc0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcW" role="lcghm">
            <property role="lacIc" value="	}{" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAcY" role="lcghm">
            <property role="lacIc" value="		Alias: (Alias)(u)," />
          </node>
          <node concept="l8MVK" id="7tgPrsAcZ" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.Sensitive) ─── -->
        <node concept="lc7rE" id="7tgPrsAc5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc1" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAc2" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAc3" role="lcghm">
            <property role="lacIc" value=": &quot;[REDACTED]&quot;," />
          </node>
          <node concept="l8MVK" id="7tgPrsAc4" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAdb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc6" role="lcghm">
            <property role="lacIc" value="	})" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAc8" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc9" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAda" role="lcghm" />
        </node>
        <!-- END IF -->
        <node concept="3clFbH" id="7tgPrsAdc" role="3cqZAp" />
        <node concept="3clFbH" id="7tgPrsAdd" role="3cqZAp" />
        <!-- // Create struct -->
        <node concept="lc7rE" id="7tgPrsAdi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAde" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAdf" role="lcghm">
            <property role="lacIc" value="▶schema.createStructName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdg" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdh" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <node concept="lc7rE" id="7tgPrsAdr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdj" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdk" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdl" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAdm" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdn" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdo" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdp" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdq" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAdv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAds" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdt" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAdu" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAdw" role="3cqZAp" />
        <!-- // ~~~~ Join schema structs ~~~~ -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <node concept="lc7rE" id="7tgPrsAdB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdx" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAdy" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdz" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdA" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAdK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdC" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdD" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdE" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdF" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdG" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdH" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdI" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdJ" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAdV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdL" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdM" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdN" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAdO" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdP" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdQ" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdR" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdS" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAdT" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdU" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAdZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdW" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdX" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAdY" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAd0" role="3cqZAp" />
        <!-- VAR: int refCount = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- ─── IF: if (refCount == 1) ─── -->
        <node concept="lc7rE" id="7tgPrsAee" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd1" role="lcghm">
            <property role="lacIc" value="type Assign" />
          </node>
          <node concept="la8eA" id="7tgPrsAd2" role="lcghm">
            <property role="lacIc" value="▶fr.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAd3" role="lcghm">
            <property role="lacIc" value="Body struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAd4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAd5" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAd6" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAd7" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAd8" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAd9" role="lcghm">
            <property role="lacIc" value="&quot; binding:&quot;required&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAea" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAeb" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAec" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAed" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- ASSIGN: refCount = refCount + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAef" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — regular -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAeg" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAeo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeh" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAei" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAej" role="lcghm">
            <property role="lacIc" value="// Repositories" />
          </node>
          <node concept="l8MVK" id="7tgPrsAek" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAel" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAem" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAen" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAep" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="3clFbH" id="7tgPrsAeq" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAew" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAer" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAes" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAet" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeu" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAev" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAex" role="3cqZAp" />
        <!-- // Create -->
        <node concept="lc7rE" id="7tgPrsAeJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAey" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAez" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAeA" role="lcghm">
            <property role="lacIc" value=") Create(u *" />
          </node>
          <node concept="la8eA" id="7tgPrsAeB" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAeC" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAeE" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAeG" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
          <node concept="la8eA" id="7tgPrsAeH" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAeI" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <!-- VAR: int idx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <!-- IF: if (idx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAeL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeK" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAeN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeM" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- ASSIGN: idx = idx + 1; -->
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAeR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeO" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAeQ" role="lcghm">
            <property role="lacIc" value="		 VALUES (" />
          </node>
        </node>
        <!-- ═══ FOR: for (int i = 1; i &lt;= idx; i++) ═══ -->
        <!-- IF: if (i &gt; 1) -->
        <node concept="lc7rE" id="7tgPrsAeT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeS" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAeW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeU" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
          <node concept="la8eA" id="7tgPrsAeV" role="lcghm">
            <property role="lacIc" value="▶i◀" />
          </node>
        </node>
        <!-- END FOR -->
        <node concept="lc7rE" id="7tgPrsAe0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeX" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAeZ" role="lcghm">
            <property role="lacIc" value="		 RETURNING id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType == DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAe3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe1" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAe2" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAe6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe4" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
          <node concept="l8MVK" id="7tgPrsAe5" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAe7" role="3cqZAp" />
        <!-- // non timestapm -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType != DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAfc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe8" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
          <node concept="la8eA" id="7tgPrsAe9" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAfa" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAfb" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAfd" role="3cqZAp" />
        <!-- // scanning -->
        <node concept="lc7rE" id="7tgPrsAff" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfe" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType == DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAfi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfg" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
          <node concept="la8eA" id="7tgPrsAfh" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAfo" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfj" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfl" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfm" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAfn" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfp" role="3cqZAp" />
        <!-- // GetByID -->
        <node concept="lc7rE" id="7tgPrsAfD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfq" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAfr" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAfs" role="lcghm">
            <property role="lacIc" value=") GetByID(id int64) (*" />
          </node>
          <node concept="la8eA" id="7tgPrsAft" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAfu" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfw" role="lcghm">
            <property role="lacIc" value="	u := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAfx" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAfy" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfA" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfC" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAfG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfE" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAfF" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAfN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfI" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAfJ" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAfK" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id," />
          </node>
          <node concept="l8MVK" id="7tgPrsAfL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfM" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAfQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfO" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
          <node concept="la8eA" id="7tgPrsAfP" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAfY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfR" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfS" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfT" role="lcghm">
            <property role="lacIc" value="	return u, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfU" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfV" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfW" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAfX" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfZ" role="3cqZAp" />
        <!-- // List -->
        <node concept="lc7rE" id="7tgPrsAf9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf0" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAf1" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAf2" role="lcghm">
            <property role="lacIc" value=") List() ([]" />
          </node>
          <node concept="la8eA" id="7tgPrsAf3" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAf4" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAf6" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAf8" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAgc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAga" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAgb" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAgz" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAge" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAgf" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAgg" role="lcghm">
            <property role="lacIc" value=" ORDER BY id`)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgh" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgi" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgk" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgl" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgm" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgo" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgp" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgq" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
          <node concept="la8eA" id="7tgPrsAgr" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgs" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgt" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgu" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgv" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
          <node concept="la8eA" id="7tgPrsAgw" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgx" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgy" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAgC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgA" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
          <node concept="la8eA" id="7tgPrsAgB" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAgS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgD" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgF" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgH" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgJ" role="lcghm">
            <property role="lacIc" value="		items = append(items, u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgL" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgN" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgO" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAgP" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgQ" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAgR" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgT" role="3cqZAp" />
        <!-- // Update -->
        <node concept="lc7rE" id="7tgPrsAg5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgU" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAgV" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAgW" role="lcghm">
            <property role="lacIc" value=") Update(u *" />
          </node>
          <node concept="la8eA" id="7tgPrsAgX" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAgY" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAg0" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAg1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAg2" role="lcghm">
            <property role="lacIc" value="		`UPDATE " />
          </node>
          <node concept="la8eA" id="7tgPrsAg3" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAg4" role="lcghm">
            <property role="lacIc" value=" SET " />
          </node>
        </node>
        <!-- VAR: int uidx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <!-- IF: if (uidx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAg7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg6" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: uidx = uidx + 1; -->
        <node concept="lc7rE" id="7tgPrsAhb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg8" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAg9" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
          <node concept="la8eA" id="7tgPrsAha" role="lcghm">
            <property role="lacIc" value="▶uidx◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- VAR: int nextParam = uidx + 1; -->
        <node concept="lc7rE" id="7tgPrsAhj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhc" role="lcghm">
            <property role="lacIc" value=", updated_at = NOW()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhd" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAhe" role="lcghm">
            <property role="lacIc" value="		 WHERE id = $" />
          </node>
          <node concept="la8eA" id="7tgPrsAhf" role="lcghm">
            <property role="lacIc" value="▶nextParam◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAhh" role="lcghm">
            <property role="lacIc" value="		 RETURNING updated_at`," />
          </node>
          <node concept="l8MVK" id="7tgPrsAhi" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType != DataType.timestamp) -->
        <node concept="lc7rE" id="7tgPrsAho" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhk" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
          <node concept="la8eA" id="7tgPrsAhl" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAhm" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAhn" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAhw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhp" role="lcghm">
            <property role="lacIc" value="		u.ID," />
          </node>
          <node concept="l8MVK" id="7tgPrsAhq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAhr" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.UpdatedAt)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhs" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAht" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhu" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAhv" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAhx" role="3cqZAp" />
        <!-- // Delete -->
        <node concept="lc7rE" id="7tgPrsAhL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhy" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAhz" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAhA" role="lcghm">
            <property role="lacIc" value=") Delete(id int64) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAhC" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAhD" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAhE" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAhG" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAhI" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhJ" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAhK" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAhM" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAhN" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="3clFbH" id="7tgPrsAhO" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAhU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhP" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAhQ" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAhR" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhS" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAhT" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAhV" role="3cqZAp" />
        <!-- // Assign -->
        <node concept="lc7rE" id="7tgPrsAhZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhW" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAhX" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAhY" role="lcghm">
            <property role="lacIc" value=") Assign(" />
          </node>
        </node>
        <!-- VAR: int argIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (argIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAh1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh0" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAh4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh2" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAh3" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <!-- ASSIGN: argIdx = argIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- // query -->
        <node concept="lc7rE" id="7tgPrsAii" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh5" role="lcghm">
            <property role="lacIc" value=") (*" />
          </node>
          <node concept="la8eA" id="7tgPrsAh6" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAh7" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAh8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAh9" role="lcghm">
            <property role="lacIc" value="	ur := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAia" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAib" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAic" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAid" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAie" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAif" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
          <node concept="la8eA" id="7tgPrsAig" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAih" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <!-- VAR: int fkIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (fkIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAik" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAij" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAim" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAil" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- ASSIGN: fkIdx = fkIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAio" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAin" role="lcghm">
            <property role="lacIc" value=") VALUES (" />
          </node>
        </node>
        <!-- ═══ FOR: for (int i = 1; i &lt;= fkIdx; i++) ═══ -->
        <!-- IF: if (i &gt; 1) -->
        <node concept="lc7rE" id="7tgPrsAiq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAip" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAit" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAir" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
          <node concept="la8eA" id="7tgPrsAis" role="lcghm">
            <property role="lacIc" value="▶i◀" />
          </node>
        </node>
        <!-- END FOR -->
        <node concept="lc7rE" id="7tgPrsAix" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiu" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAiw" role="lcghm">
            <property role="lacIc" value="		 ON CONFLICT (" />
          </node>
        </node>
        <!-- VAR: int ckIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (ckIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAiz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiy" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAiB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiA" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- ASSIGN: ckIdx = ckIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAiF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiC" role="lcghm">
            <property role="lacIc" value=") DO NOTHING" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAiE" role="lcghm">
            <property role="lacIc" value="		 RETURNING " />
          </node>
        </node>
        <!-- VAR: int retIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- IF: if (retIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAiH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiG" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAiJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiI" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAiL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiK" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: retIdx = retIdx + 1; -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAiO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiM" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
          <node concept="l8MVK" id="7tgPrsAiN" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAiT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiP" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAiQ" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAiR" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAiS" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAiV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiU" role="lcghm">
            <property role="lacIc" value="	).Scan(" />
          </node>
        </node>
        <!-- VAR: int scanIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- IF: if (scanIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAiX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiW" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAi0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiY" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
          <node concept="la8eA" id="7tgPrsAiZ" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="7tgPrsAi3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi1" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
          <node concept="la8eA" id="7tgPrsAi2" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: scanIdx = scanIdx + 1; -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAjb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi4" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAi5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAi6" role="lcghm">
            <property role="lacIc" value="	return ur, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAi7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAi8" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAi9" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAja" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAjc" role="3cqZAp" />
        <!-- // Remove -->
        <node concept="lc7rE" id="7tgPrsAjg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjd" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAje" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjf" role="lcghm">
            <property role="lacIc" value=") Remove(" />
          </node>
        </node>
        <!-- VAR: int rmIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (rmIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAji" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjh" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="7tgPrsAjl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjj" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjk" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <!-- ASSIGN: rmIdx = rmIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAjr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjm" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjo" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAjp" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjq" role="lcghm">
            <property role="lacIc" value=" WHERE " />
          </node>
        </node>
        <!-- VAR: int wIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (wIdx &gt; 0) -->
        <node concept="lc7rE" id="7tgPrsAjt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjs" role="lcghm">
            <property role="lacIc" value=" AND " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: wIdx = wIdx + 1; -->
        <node concept="lc7rE" id="7tgPrsAjx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAju" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjv" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
          <node concept="la8eA" id="7tgPrsAjw" role="lcghm">
            <property role="lacIc" value="▶wIdx◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAjz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjy" role="lcghm">
            <property role="lacIc" value="`" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="7tgPrsAjC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjA" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAjB" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAjK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjD" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjF" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjH" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjI" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAjJ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAjL" role="3cqZAp" />
        <!-- // Cross-queries -->
        <!-- ═══ FOREACH: foreach refA in schema.Fields ═══ -->
        <!-- ─── IF: if (refA.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA; -->
        <!-- ═══ FOREACH: foreach refB in schema.Fields ═══ -->
        <!-- ─── IF: if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB; -->
        <!-- VAR: string ts = frB.target_schema.structName(); -->
        <!-- VAR: string tt = frB.target_schema.name; -->
        <node concept="3clFbH" id="7tgPrsAjM" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAj2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjN" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
          <node concept="la8eA" id="7tgPrsAjO" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjP" role="lcghm">
            <property role="lacIc" value=") Get" />
          </node>
          <node concept="la8eA" id="7tgPrsAjQ" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjR" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
          <node concept="la8eA" id="7tgPrsAjS" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjT" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAjU" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjV" role="lcghm">
            <property role="lacIc" value=" int64) ([]" />
          </node>
          <node concept="la8eA" id="7tgPrsAjW" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAjX" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAjZ" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAj0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAj1" role="lcghm">
            <property role="lacIc" value="		`SELECT t.id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach tf in frB.target_schema.Fields ═══ -->
        <!-- ─── IF: if (tf.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) tf; -->
        <node concept="lc7rE" id="7tgPrsAj5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj3" role="lcghm">
            <property role="lacIc" value=", t." />
          </node>
          <node concept="la8eA" id="7tgPrsAj4" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAkI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAj6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAj7" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
          <node concept="la8eA" id="7tgPrsAj8" role="lcghm">
            <property role="lacIc" value="▶tt◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAj9" role="lcghm">
            <property role="lacIc" value=" t" />
          </node>
          <node concept="l8MVK" id="7tgPrsAka" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkb" role="lcghm">
            <property role="lacIc" value="		 INNER JOIN " />
          </node>
          <node concept="la8eA" id="7tgPrsAkc" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAkd" role="lcghm">
            <property role="lacIc" value=" j ON j." />
          </node>
          <node concept="la8eA" id="7tgPrsAke" role="lcghm">
            <property role="lacIc" value="▶frB.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAkf" role="lcghm">
            <property role="lacIc" value=" = t.id" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkh" role="lcghm">
            <property role="lacIc" value="		 WHERE j." />
          </node>
          <node concept="la8eA" id="7tgPrsAki" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAkj" role="lcghm">
            <property role="lacIc" value=" = $1" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkl" role="lcghm">
            <property role="lacIc" value="		 ORDER BY t.id`, " />
          </node>
          <node concept="la8eA" id="7tgPrsAkm" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAkn" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAko" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkp" role="lcghm">
            <property role="lacIc" value="	)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkr" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAks" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkt" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAku" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkv" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkw" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkx" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAky" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkz" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
          <node concept="la8eA" id="7tgPrsAkA" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkC" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkE" role="lcghm">
            <property role="lacIc" value="		var item " />
          </node>
          <node concept="la8eA" id="7tgPrsAkF" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkH" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;item.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach tf in frB.target_schema.Fields ═══ -->
        <!-- ─── IF: if (tf.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) tf; -->
        <node concept="lc7rE" id="7tgPrsAkL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkJ" role="lcghm">
            <property role="lacIc" value=", &amp;item." />
          </node>
          <node concept="la8eA" id="7tgPrsAkK" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="7tgPrsAk1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkM" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkO" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkQ" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkS" role="lcghm">
            <property role="lacIc" value="		items = append(items, item)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkU" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkW" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAkY" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkZ" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAk0" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAk2" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — regular schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAk3" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <node concept="3clFbH" id="7tgPrsAk4" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAld" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk5" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAk6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAk7" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
          <node concept="la8eA" id="7tgPrsAk8" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAk9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAla" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlb" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAlc" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAle" role="3cqZAp" />
        <!-- // Create handler -->
        <node concept="lc7rE" id="7tgPrsAlR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlf" role="lcghm">
            <property role="lacIc" value="func handleCreate" />
          </node>
          <node concept="la8eA" id="7tgPrsAlg" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAlh" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAli" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAlj" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAll" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAln" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
          <node concept="la8eA" id="7tgPrsAlo" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlp" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlq" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlr" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAls" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlt" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlu" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlw" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlx" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAly" role="lcghm">
            <property role="lacIc" value="		if err := repo.Create(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlA" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlC" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlE" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlG" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlI" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlJ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlK" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlM" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlO" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlP" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAlQ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAlS" role="3cqZAp" />
        <!-- // Get handler -->
        <node concept="lc7rE" id="7tgPrsAmw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlT" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAlU" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAlV" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAlW" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAlX" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAlZ" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAl1" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAl3" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAl5" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAl7" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAl9" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAma" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmb" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmd" role="lcghm">
            <property role="lacIc" value="		u, err := repo.GetByID(id)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAme" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmf" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmh" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusNotFound)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmi" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmj" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAml" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmn" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmo" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmp" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmr" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAms" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmt" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmu" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAmv" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAmx" role="3cqZAp" />
        <!-- // List handler -->
        <node concept="lc7rE" id="7tgPrsAmZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmy" role="lcghm">
            <property role="lacIc" value="func handleList" />
          </node>
          <node concept="la8eA" id="7tgPrsAmz" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAmA" role="lcghm">
            <property role="lacIc" value="s(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmB" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAmC" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmE" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmG" role="lcghm">
            <property role="lacIc" value="		items, err := repo.List()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmI" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmJ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmK" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmM" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmO" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmQ" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmS" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmU" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAmW" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmX" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAmY" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAm0" role="3cqZAp" />
        <!-- // Update handler -->
        <node concept="lc7rE" id="7tgPrsAnP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm1" role="lcghm">
            <property role="lacIc" value="func handleUpdate" />
          </node>
          <node concept="la8eA" id="7tgPrsAm2" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAm3" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAm4" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAm5" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAm6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAm7" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAm8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAm9" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAna" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnb" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnd" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAne" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnf" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAng" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnh" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAni" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnj" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnl" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
          <node concept="la8eA" id="7tgPrsAnm" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAno" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnp" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnq" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnr" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAns" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnt" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnu" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnw" role="lcghm">
            <property role="lacIc" value="		u.ID = id" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnx" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAny" role="lcghm">
            <property role="lacIc" value="		if err := repo.Update(&amp;u); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnA" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnB" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnC" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnD" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnE" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnF" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnG" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnH" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnI" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnJ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnK" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnM" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnN" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAnO" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAnQ" role="3cqZAp" />
        <!-- // Delete handler -->
        <node concept="lc7rE" id="7tgPrsAoq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnR" role="lcghm">
            <property role="lacIc" value="func handleDelete" />
          </node>
          <node concept="la8eA" id="7tgPrsAnS" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAnT" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAnU" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAnV" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnW" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnX" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAnZ" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAn1" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAn3" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAn5" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAn7" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAn9" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoa" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAob" role="lcghm">
            <property role="lacIc" value="		if err := repo.Delete(id); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAod" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoe" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAof" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAog" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoh" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoi" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoj" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAok" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAol" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAom" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAon" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoo" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAop" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAor" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAos" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <node concept="3clFbH" id="7tgPrsAot" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAoD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAou" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAov" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAow" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
          <node concept="la8eA" id="7tgPrsAox" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAoy" role="lcghm">
            <property role="lacIc" value=" (assignments)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoA" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoB" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAoC" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAoE" role="3cqZAp" />
        <!-- VAR: node&lt;FieldRefrence&gt; firstRef = null; -->
        <!-- VAR: node&lt;FieldRefrence&gt; secondRef = null; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (firstRef == null) { firstRef = fr; } -->
        <!-- IF: else if (secondRef == null) { secondRef = fr; } -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAoF" role="3cqZAp" />
        <!-- // Assign handler -->
        <node concept="lc7rE" id="7tgPrsApB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoG" role="lcghm">
            <property role="lacIc" value="func handleAssign" />
          </node>
          <node concept="la8eA" id="7tgPrsAoH" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAoI" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAoJ" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAoK" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoL" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoM" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoO" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAoP" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAoQ" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoS" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoU" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoW" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAoY" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAo0" role="lcghm">
            <property role="lacIc" value="		var body Assign" />
          </node>
          <node concept="la8eA" id="7tgPrsAo1" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAo2" role="lcghm">
            <property role="lacIc" value="Body" />
          </node>
          <node concept="l8MVK" id="7tgPrsAo3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAo4" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAo5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAo6" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAo7" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAo8" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAo9" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApa" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsApb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApc" role="lcghm">
            <property role="lacIc" value="		ur, err := urRepo.Assign(" />
          </node>
          <node concept="la8eA" id="7tgPrsApd" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsApe" role="lcghm">
            <property role="lacIc" value=", body." />
          </node>
          <node concept="la8eA" id="7tgPrsApf" role="lcghm">
            <property role="lacIc" value="▶secondRef.pascalName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsApg" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAph" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApi" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsApj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApk" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsApl" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApm" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsApn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApo" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsApp" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApq" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsApr" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAps" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
          <node concept="l8MVK" id="7tgPrsApt" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApu" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(ur)" />
          </node>
          <node concept="l8MVK" id="7tgPrsApv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApw" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsApx" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApy" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsApz" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsApA" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsApC" role="3cqZAp" />
        <!-- // Remove handler -->
        <node concept="lc7rE" id="7tgPrsAqu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApD" role="lcghm">
            <property role="lacIc" value="func handleRemove" />
          </node>
          <node concept="la8eA" id="7tgPrsApE" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsApF" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
          <node concept="la8eA" id="7tgPrsApG" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsApH" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsApI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApJ" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsApK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApL" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsApM" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsApN" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsApO" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApP" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsApQ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApR" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsApS" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApT" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsApU" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApV" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsApW" role="lcghm" />
          <node concept="la8eA" id="7tgPrsApX" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsApY" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsApZ" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAp0" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAp1" role="lcghm">
            <property role="lacIc" value="&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAp2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAp3" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAp4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAp5" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAp6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAp7" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAp8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAp9" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqa" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqb" role="lcghm">
            <property role="lacIc" value="		if err := urRepo.Remove(" />
          </node>
          <node concept="la8eA" id="7tgPrsAqc" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAqd" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
          <node concept="la8eA" id="7tgPrsAqe" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAqf" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqh" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqi" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqj" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqk" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAql" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqm" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqn" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqo" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqp" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqq" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqr" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqs" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAqt" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAqv" role="3cqZAp" />
        <!-- // Cross-query handlers -->
        <!-- ═══ FOREACH: foreach refA in schema.Fields ═══ -->
        <!-- ─── IF: if (refA.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA; -->
        <!-- ═══ FOREACH: foreach refB in schema.Fields ═══ -->
        <!-- ─── IF: if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB; -->
        <node concept="lc7rE" id="7tgPrsArc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAqw" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAqx" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAqy" role="lcghm">
            <property role="lacIc" value="▶frB.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAqz" role="lcghm">
            <property role="lacIc" value="s(urRepo *" />
          </node>
          <node concept="la8eA" id="7tgPrsAqA" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAqB" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqD" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqF" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqH" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqJ" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqL" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqM" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqN" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqO" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqP" role="lcghm">
            <property role="lacIc" value="		items, err := urRepo.Get" />
          </node>
          <node concept="la8eA" id="7tgPrsAqQ" role="lcghm">
            <property role="lacIc" value="▶frB.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAqR" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
          <node concept="la8eA" id="7tgPrsAqS" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAqT" role="lcghm">
            <property role="lacIc" value="(id)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqU" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqV" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqW" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqX" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAqY" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAqZ" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAq0" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAq1" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAq2" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAq3" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAq4" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAq5" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAq6" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAq7" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAq8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAq9" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAra" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsArb" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsArd" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Main -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="7tgPrsAre" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAsi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsArf" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsArg" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArh" role="lcghm">
            <property role="lacIc" value="// Main" />
          </node>
          <node concept="l8MVK" id="7tgPrsAri" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArj" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsArk" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsArl" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArm" role="lcghm">
            <property role="lacIc" value="func main() {" />
          </node>
          <node concept="l8MVK" id="7tgPrsArn" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAro" role="lcghm">
            <property role="lacIc" value="	dbURL := os.Getenv(&quot;DATABASE_URL&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsArp" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArq" role="lcghm">
            <property role="lacIc" value="	if dbURL == &quot;&quot; {" />
          </node>
          <node concept="l8MVK" id="7tgPrsArr" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArs" role="lcghm">
            <property role="lacIc" value="		dbURL = &quot;postgres://" />
          </node>
          <node concept="la8eA" id="7tgPrsArt" role="lcghm">
            <property role="lacIc" value="▶infra.dbUser◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAru" role="lcghm">
            <property role="lacIc" value=":" />
          </node>
          <node concept="la8eA" id="7tgPrsArv" role="lcghm">
            <property role="lacIc" value="▶infra.dbPass◀" />
          </node>
          <node concept="la8eA" id="7tgPrsArw" role="lcghm">
            <property role="lacIc" value="@localhost:5432/" />
          </node>
          <node concept="la8eA" id="7tgPrsArx" role="lcghm">
            <property role="lacIc" value="▶infra.dbName◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAry" role="lcghm">
            <property role="lacIc" value="?sslmode=disable&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsArz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArA" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsArB" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsArC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArD" role="lcghm">
            <property role="lacIc" value="	db, err := sql.Open(&quot;postgres&quot;, dbURL)" />
          </node>
          <node concept="l8MVK" id="7tgPrsArE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArF" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsArG" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArH" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsArI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArJ" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsArK" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArL" role="lcghm">
            <property role="lacIc" value="	defer db.Close()" />
          </node>
          <node concept="l8MVK" id="7tgPrsArM" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsArN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArO" role="lcghm">
            <property role="lacIc" value="	for i := 0; i &lt; 5; i++ {" />
          </node>
          <node concept="l8MVK" id="7tgPrsArP" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArQ" role="lcghm">
            <property role="lacIc" value="		if err = db.Ping(); err == nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsArR" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArS" role="lcghm">
            <property role="lacIc" value="			break" />
          </node>
          <node concept="l8MVK" id="7tgPrsArT" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArU" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
          <node concept="l8MVK" id="7tgPrsArV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArW" role="lcghm">
            <property role="lacIc" value="		log.Printf(&quot;DB not ready, retrying... (%d/5)&quot;, i+1)" />
          </node>
          <node concept="l8MVK" id="7tgPrsArX" role="lcghm" />
          <node concept="la8eA" id="7tgPrsArY" role="lcghm">
            <property role="lacIc" value="		time.Sleep(2 * time.Second)" />
          </node>
          <node concept="l8MVK" id="7tgPrsArZ" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAr0" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAr1" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAr2" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAr3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAr4" role="lcghm">
            <property role="lacIc" value="		log.Fatal(&quot;DB connection failed:&quot;, err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAr5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAr6" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAr7" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAr8" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAr9" role="lcghm">
            <property role="lacIc" value="	if _, err := db.Exec(migrationSQL); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsa" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsb" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsc" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsd" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAse" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsf" role="lcghm">
            <property role="lacIc" value="	log.Println(&quot;Migration complete&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsg" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAsh" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAsj" role="3cqZAp" />
        <!-- // Repo instantiation -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <node concept="lc7rE" id="7tgPrsAsq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsk" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAsl" role="lcghm">
            <property role="lacIc" value="▶schema.singularName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAsm" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAsn" role="lcghm">
            <property role="lacIc" value="▶schema.repoName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAso" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsp" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <node concept="lc7rE" id="7tgPrsAsx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsr" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAss" role="lcghm">
            <property role="lacIc" value="▶schema.singularName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAst" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAsu" role="lcghm">
            <property role="lacIc" value="▶schema.repoName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAsv" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsw" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAsy" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAsI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAsz" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsA" role="lcghm">
            <property role="lacIc" value="	mux := http.NewServeMux()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsB" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAsC" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsD" role="lcghm">
            <property role="lacIc" value="	// Swagger UI" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsE" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsF" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /swagger/*&quot;, httpSwagger.WrapHandler)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsG" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAsH" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAsJ" role="3cqZAp" />
        <!-- // Regular routes -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string vr = schema.singularName() + "Repo"; -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="lc7rE" id="7tgPrsAtt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAsK" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
          <node concept="la8eA" id="7tgPrsAsL" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAsM" role="lcghm">
            <property role="lacIc" value="s" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsN" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsO" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
          <node concept="la8eA" id="7tgPrsAsP" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAsQ" role="lcghm">
            <property role="lacIc" value="&quot;, handleCreate" />
          </node>
          <node concept="la8eA" id="7tgPrsAsR" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAsS" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAsT" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAsU" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAsV" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAsW" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsAsX" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAsY" role="lcghm">
            <property role="lacIc" value="&quot;, handleList" />
          </node>
          <node concept="la8eA" id="7tgPrsAsZ" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAs0" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
          <node concept="la8eA" id="7tgPrsAs1" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAs2" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAs3" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAs4" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsAs5" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAs6" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAs7" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAs8" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAs9" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAta" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtb" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtc" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;PUT /" />
          </node>
          <node concept="la8eA" id="7tgPrsAtd" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAte" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleUpdate" />
          </node>
          <node concept="la8eA" id="7tgPrsAtf" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtg" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAth" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAti" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtj" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtk" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
          <node concept="la8eA" id="7tgPrsAtl" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtm" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleDelete" />
          </node>
          <node concept="la8eA" id="7tgPrsAtn" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAto" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAtp" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtq" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtr" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAts" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAtu" role="3cqZAp" />
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
        <node concept="lc7rE" id="7tgPrsAui" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAtv" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
          <node concept="la8eA" id="7tgPrsAtw" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtx" role="lcghm">
            <property role="lacIc" value=" assignments" />
          </node>
          <node concept="l8MVK" id="7tgPrsAty" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtz" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
          <node concept="la8eA" id="7tgPrsAtA" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtB" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsAtC" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtD" role="lcghm">
            <property role="lacIc" value="&quot;, handleAssign" />
          </node>
          <node concept="la8eA" id="7tgPrsAtE" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtF" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAtG" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtH" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtI" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtJ" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
          <node concept="la8eA" id="7tgPrsAtK" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtL" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsAtM" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtN" role="lcghm">
            <property role="lacIc" value="/{" />
          </node>
          <node concept="la8eA" id="7tgPrsAtO" role="lcghm">
            <property role="lacIc" value="▶sRef.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtP" role="lcghm">
            <property role="lacIc" value="}&quot;, handleRemove" />
          </node>
          <node concept="la8eA" id="7tgPrsAtQ" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtR" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
          <node concept="la8eA" id="7tgPrsAtS" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtT" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAtU" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAtV" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsAtW" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtX" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsAtY" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAtZ" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAt0" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAt1" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAt2" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
          <node concept="la8eA" id="7tgPrsAt3" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAt4" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAt5" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAt6" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
          <node concept="la8eA" id="7tgPrsAt7" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAt8" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
          <node concept="la8eA" id="7tgPrsAt9" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAua" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
          <node concept="la8eA" id="7tgPrsAub" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAuc" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.structName()◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAud" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
          <node concept="la8eA" id="7tgPrsAue" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAuf" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAug" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAuh" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="7tgPrsAuj" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAuy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAuk" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Serving on :" />
          </node>
          <node concept="la8eA" id="7tgPrsAul" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAum" role="lcghm">
            <property role="lacIc" value="&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAun" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuo" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Swagger UI: http://localhost:" />
          </node>
          <node concept="la8eA" id="7tgPrsAup" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAuq" role="lcghm">
            <property role="lacIc" value="/swagger/index.html&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAur" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAus" role="lcghm">
            <property role="lacIc" value="	log.Fatal(http.ListenAndServe(&quot;:" />
          </node>
          <node concept="la8eA" id="7tgPrsAut" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
          <node concept="la8eA" id="7tgPrsAuu" role="lcghm">
            <property role="lacIc" value="&quot;, mux))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAuv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAuw" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAux" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAuz" role="3cqZAp" />
        <!-- END ? -->
        <!-- END ? -->
      </node>
    </node>
  </node>
</model>
