<?xml version="1.0" encoding="UTF-8"?>
<model ref="r:f10fb75e-9e61-4ecf-8328-e03d2ea8b365(UserManagement.textGen)">
  <persistence version="9" />
  <languages>
    <use id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen" version="1" />
    <devkit ref="fa73d85a-ac7f-447b-846c-fcdc41caa600(jetbrains.mps.devkit.aspect.textgen)" />
  </languages>
  <imports>
    <import index="s1" ref="r:40b69eb5-3c02-4fef-92ca-f6d48dc1801a(throwed.structure)" implicit="true" />
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
  <node concept="WtQ9Q" id="7tgPrsApf">
    <ref role="WuzLi" to="s1:TODO-CONCEPT-ID" resolve="Entity" />
    <node concept="29tfMY" id="7tgPrsApi" role="29tGrW">
      <node concept="3clFbS" id="7tgPrsApj" role="2VODD2">
        <node concept="3cpWs6" id="7tgPrsApk" role="3cqZAp">
          <node concept="Xl_RD" id="7tgPrsApl" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="7tgPrsApg" role="11c4hB">
      <node concept="3clFbS" id="7tgPrsAph" role="2VODD2">
        <node concept="3clFbH" id="7tgPrsAa5" role="3cqZAp" />
        <!-- // === Top-level variables === -->
        <node concept="lc7rE" id="7tgPrsAa6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa7" role="lcghm">
            <property role="lacIc" value="{???-string name = node.name;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAa8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa9" role="lcghm">
            <property role="lacIc" value="{???-string nameLower = node.name.toLowerCase();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAba" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbb" role="lcghm">
            <property role="lacIc" value="{???-string pkField = node.primaryKeyField().name;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAbc" role="3cqZAp" />
        <!-- // === Package + Imports === -->
        <node concept="lc7rE" id="7tgPrsAbf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbd" role="lcghm">
            <property role="lacIc" value="package main" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbe" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbh" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbi" role="lcghm">
            <property role="lacIc" value="import (" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbl" role="lcghm">
            <property role="lacIc" value="	&quot;context&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbo" role="lcghm">
            <property role="lacIc" value="	&quot;encoding/json&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbr" role="lcghm">
            <property role="lacIc" value="	&quot;fmt&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbs" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbu" role="lcghm">
            <property role="lacIc" value="	&quot;log&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbx" role="lcghm">
            <property role="lacIc" value="	&quot;time&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAby" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbA" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbC" role="lcghm">
            <property role="lacIc" value="	&quot;github.com/nats-io/nats.go&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbG" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbH" role="lcghm">
            <property role="lacIc" value="	&quot;dev.azure.com/Motadata/NextGen/motadata-go-sdk/events&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbI" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbK" role="lcghm">
            <property role="lacIc" value="	&quot;dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbN" role="lcghm">
            <property role="lacIc" value="	&quot;dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbQ" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbT" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAbV" role="3cqZAp" />
        <!-- // === Domain Struct === -->
        <node concept="lc7rE" id="7tgPrsAb0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbW" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAbX" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbY" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAb1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb2" role="lcghm">
            <property role="lacIc" value="{???-foreach field in node.fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb3" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAb4" role="lcghm">
            <property role="lacIc" value="{???-field.name.capitalize()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb5" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAb6" role="lcghm">
            <property role="lacIc" value="{???-field.goType()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb7" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAb8" role="lcghm">
            <property role="lacIc" value="{???-field.jsonName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb9" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAca" role="lcghm">
            <property role="lacIc" value="{???-field.dbName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcb" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAce" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcf" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAci" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcg" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAch" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAck" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcj" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAcl" role="3cqZAp" />
        <!-- // === Event / Request Types === -->
        <node concept="lc7rE" id="7tgPrsAcm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcn" role="lcghm">
            <property role="lacIc" value="{???-foreach op in node.operations {}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAco" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcq" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:create) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcr" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAcs" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAct" role="lcghm">
            <property role="lacIc" value="CreatedEvent struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcw" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcx" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcy" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAcz" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcA" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcB" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcC" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcF" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcI" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcM" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcO" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAcP" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcR" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:update) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcS" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAcT" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcU" role="lcghm">
            <property role="lacIc" value="UpdatedEvent struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcX" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcY" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcZ" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
          <node concept="la8eA" id="7tgPrsAc0" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAc1" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAc2" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAc3" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc6" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc9" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAda" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAde" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdf" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdg" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdi" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:delete) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdj" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAdk" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdl" role="lcghm">
            <property role="lacIc" value="DeletedEvent struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdo" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdp" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdq" role="lcghm">
            <property role="lacIc" value="ID string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdr" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAds" role="lcghm">
            <property role="lacIc" value="_id&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdt" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdv" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdy" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdC" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdE" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdF" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdH" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdI" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAdJ" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdK" role="lcghm">
            <property role="lacIc" value="ListRequest struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdN" role="lcghm">
            <property role="lacIc" value="	Limit     int       `json:&quot;limit&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdQ" role="lcghm">
            <property role="lacIc" value="	Offset    int       `json:&quot;offset&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdT" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdW" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd0" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd2" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAd3" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAd4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd5" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:get) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAea" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd6" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAd7" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAd8" role="lcghm">
            <property role="lacIc" value="GetRequest struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAd9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeb" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAec" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAed" role="lcghm">
            <property role="lacIc" value="ID string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAee" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAef" role="lcghm">
            <property role="lacIc" value="_id&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAek" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAei" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAej" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAen" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAel" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAem" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAep" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAer" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAes" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAet" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeu" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAev" role="3cqZAp" />
        <!-- // === Handler Struct + Constructor === -->
        <node concept="lc7rE" id="7tgPrsAeA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAew" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAex" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAey" role="lcghm">
            <property role="lacIc" value="Handler struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAez" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeB" role="lcghm">
            <property role="lacIc" value="	publisher     *events.Publisher" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeE" role="lcghm">
            <property role="lacIc" value="	subjectPrefix string" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeH" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeI" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeL" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAeK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeM" role="lcghm">
            <property role="lacIc" value="func New" />
          </node>
          <node concept="la8eA" id="7tgPrsAeN" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeO" role="lcghm">
            <property role="lacIc" value="Handler(pub *events.Publisher, subjectPrefix string) *" />
          </node>
          <node concept="la8eA" id="7tgPrsAeP" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeQ" role="lcghm">
            <property role="lacIc" value="Handler {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeT" role="lcghm">
            <property role="lacIc" value="	return &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAeU" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeV" role="lcghm">
            <property role="lacIc" value="Handler{" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAe0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeY" role="lcghm">
            <property role="lacIc" value="		publisher:     pub," />
          </node>
          <node concept="l8MVK" id="7tgPrsAeZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAe3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe1" role="lcghm">
            <property role="lacIc" value="		subjectPrefix: subjectPrefix," />
          </node>
          <node concept="l8MVK" id="7tgPrsAe2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAe6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe4" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAe9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe7" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfb" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfa" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfc" role="3cqZAp" />
        <!-- // === Handler Methods === -->
        <node concept="lc7rE" id="7tgPrsAfd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfe" role="lcghm">
            <property role="lacIc" value="{???-foreach op in node.operations {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAff" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfg" role="lcghm">
            <property role="lacIc" value="{???-string opName = op.capitalizedName();}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfi" role="lcghm">
            <property role="lacIc" value="{???-string opKind = op.entityOperation.name;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAfj" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAfq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfk" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAfl" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfm" role="lcghm">
            <property role="lacIc" value="Handler) Handle" />
          </node>
          <node concept="la8eA" id="7tgPrsAfn" role="lcghm">
            <property role="lacIc" value="{???-opName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfo" role="lcghm">
            <property role="lacIc" value="(req core.Request) {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAft" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfr" role="lcghm">
            <property role="lacIc" value="	ctx := req.Context()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfs" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfu" role="lcghm">
            <property role="lacIc" value="	ctx, span := tracer.StartConsumer(ctx, &quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAfv" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfw" role="lcghm">
            <property role="lacIc" value=".Handle" />
          </node>
          <node concept="la8eA" id="7tgPrsAfx" role="lcghm">
            <property role="lacIc" value="{???-opName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfy" role="lcghm">
            <property role="lacIc" value="&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfB" role="lcghm">
            <property role="lacIc" value="	defer span.End()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfE" role="lcghm">
            <property role="lacIc" value="	ctx = core.InjectContext(ctx, req.Headers())" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfI" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfH" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfJ" role="3cqZAp" />
        <!-- // ~~- Unmarshal into correct event type ~~- -->
        <node concept="lc7rE" id="7tgPrsAfK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfL" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:create) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfM" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAfN" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfO" role="lcghm">
            <property role="lacIc" value="CreatedEvent" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfP" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfS" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfT" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfU" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:update) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfV" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAfW" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfX" role="lcghm">
            <property role="lacIc" value="UpdatedEvent" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAf0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf1" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf3" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:delete) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAf8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf4" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAf5" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAf6" role="lcghm">
            <property role="lacIc" value="DeletedEvent" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAf9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAga" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgc" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgd" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAge" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgf" role="lcghm">
            <property role="lacIc" value="ListRequest" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgj" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgl" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:get) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgm" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAgn" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgo" role="lcghm">
            <property role="lacIc" value="GetRequest" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgs" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAgt" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAgw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgu" role="lcghm">
            <property role="lacIc" value="	if err := json.Unmarshal(req.Data(), &amp;event); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgx" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgy" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgA" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;400&quot;, &quot;invalid JSON: &quot; + err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgD" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgG" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgJ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgL" role="3cqZAp" />
        <!-- // ~~- Validation ~~- -->
        <node concept="lc7rE" id="7tgPrsAgM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgN" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:create) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgP" role="lcghm">
            <property role="lacIc" value="{???-int valIdx = 0;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgQ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgR" role="lcghm">
            <property role="lacIc" value="{???-foreach field in node.fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgT" role="lcghm">
            <property role="lacIc" value="{???-if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:auto)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:hidden)) &amp;&amp; !(field.hasAnnotation(FieldAnnotation:nullable))) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgV" role="lcghm">
            <property role="lacIc" value="{???-if (valIdx == 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgW" role="lcghm">
            <property role="lacIc" value="	if event." />
          </node>
          <node concept="la8eA" id="7tgPrsAgX" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgY" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAgZ" role="lcghm">
            <property role="lacIc" value="{???-field.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAg0" role="lcghm">
            <property role="lacIc" value=" == &quot;&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg3" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAg4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg5" role="lcghm">
            <property role="lacIc" value="{???-if (valIdx &gt; 0) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg6" role="lcghm">
            <property role="lacIc" value=" || event." />
          </node>
          <node concept="la8eA" id="7tgPrsAg7" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAg8" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAg9" role="lcghm">
            <property role="lacIc" value="{???-field.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAha" role="lcghm">
            <property role="lacIc" value=" == &quot;&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhd" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhf" role="lcghm">
            <property role="lacIc" value="{???-valIdx = valIdx + 1;}" />
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
        <node concept="lc7rE" id="7tgPrsAhm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhk" role="lcghm">
            <property role="lacIc" value=" {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhn" role="lcghm">
            <property role="lacIc" value="		err := fmt.Errorf(&quot;invalid " />
          </node>
          <node concept="la8eA" id="7tgPrsAho" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAhp" role="lcghm">
            <property role="lacIc" value=" data: missing required fields&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhs" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAht" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhv" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;400&quot;, err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhy" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhB" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhF" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAhG" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAhH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhI" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:update) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhJ" role="lcghm">
            <property role="lacIc" value="	if event." />
          </node>
          <node concept="la8eA" id="7tgPrsAhK" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAhL" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAhM" role="lcghm">
            <property role="lacIc" value="{???-pkField}" />
          </node>
          <node concept="la8eA" id="7tgPrsAhN" role="lcghm">
            <property role="lacIc" value=" == &quot;&quot; {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhQ" role="lcghm">
            <property role="lacIc" value="		err := fmt.Errorf(&quot;invalid " />
          </node>
          <node concept="la8eA" id="7tgPrsAhR" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAhS" role="lcghm">
            <property role="lacIc" value=" data: missing ID&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhV" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAh0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhY" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;400&quot;, err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAh3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh1" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAh2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAh6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh4" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAh5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAh7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAh8" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAh9" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAia" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAib" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:delete) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAig" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAic" role="lcghm">
            <property role="lacIc" value="	if event." />
          </node>
          <node concept="la8eA" id="7tgPrsAid" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAie" role="lcghm">
            <property role="lacIc" value="ID == &quot;&quot; {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAif" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAil" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAih" role="lcghm">
            <property role="lacIc" value="		err := fmt.Errorf(&quot;invalid request: missing " />
          </node>
          <node concept="la8eA" id="7tgPrsAii" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAij" role="lcghm">
            <property role="lacIc" value=" ID&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAik" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAio" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAim" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAin" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAir" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAip" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;400&quot;, err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAis" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAit" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAix" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiv" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiz" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAiA" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAiB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiC" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:get) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAiH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiD" role="lcghm">
            <property role="lacIc" value="	if event." />
          </node>
          <node concept="la8eA" id="7tgPrsAiE" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAiF" role="lcghm">
            <property role="lacIc" value="ID == &quot;&quot; {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiI" role="lcghm">
            <property role="lacIc" value="		err := fmt.Errorf(&quot;invalid request: missing " />
          </node>
          <node concept="la8eA" id="7tgPrsAiJ" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAiK" role="lcghm">
            <property role="lacIc" value=" ID&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiN" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiQ" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;400&quot;, err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiT" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAiW" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAiX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAiZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi0" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAi1" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAi2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi3" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAi6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi4" role="lcghm">
            <property role="lacIc" value="	if event.Limit &lt; 0 || event.Offset &lt; 0 {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAi5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAi9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAi7" role="lcghm">
            <property role="lacIc" value="		err := fmt.Errorf(&quot;invalid pagination parameters&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAi8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAja" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjd" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;400&quot;, err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAje" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAji" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjg" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjj" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjn" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAjo" role="3cqZAp" />
        <!-- // ~~- Span Attributes ~~- -->
        <node concept="lc7rE" id="7tgPrsAjq" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAjp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjr" role="lcghm">
            <property role="lacIc" value="	span.SetAttributes(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAjs" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAju" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjv" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:create) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjw" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAjx" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjy" role="lcghm">
            <property role="lacIc" value=".id&quot;, event." />
          </node>
          <node concept="la8eA" id="7tgPrsAjz" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjA" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAjB" role="lcghm">
            <property role="lacIc" value="{???-pkField}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjC" role="lcghm">
            <property role="lacIc" value=")," />
          </node>
          <node concept="l8MVK" id="7tgPrsAjD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjG" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjI" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:update) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjJ" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAjK" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjL" role="lcghm">
            <property role="lacIc" value=".id&quot;, event." />
          </node>
          <node concept="la8eA" id="7tgPrsAjM" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjN" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAjO" role="lcghm">
            <property role="lacIc" value="{???-pkField}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjP" role="lcghm">
            <property role="lacIc" value=")," />
          </node>
          <node concept="l8MVK" id="7tgPrsAjQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAjS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjT" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAjU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjV" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:delete) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAj2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAjW" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAjX" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAjY" role="lcghm">
            <property role="lacIc" value=".id&quot;, event." />
          </node>
          <node concept="la8eA" id="7tgPrsAjZ" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAj0" role="lcghm">
            <property role="lacIc" value="ID)," />
          </node>
          <node concept="l8MVK" id="7tgPrsAj1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAj3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj4" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAj5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj6" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:get) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAkd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAj7" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAj8" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAj9" role="lcghm">
            <property role="lacIc" value=".id&quot;, event." />
          </node>
          <node concept="la8eA" id="7tgPrsAka" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAkb" role="lcghm">
            <property role="lacIc" value="ID)," />
          </node>
          <node concept="l8MVK" id="7tgPrsAkc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAke" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkf" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAki" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkg" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;tenant.id&quot;, req.Header(core.HeaderTenantID))," />
          </node>
          <node concept="l8MVK" id="7tgPrsAkh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkj" role="lcghm">
            <property role="lacIc" value="	)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkn" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkm" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAko" role="3cqZAp" />
        <!-- // ~~- Pre-hook ~~- -->
        <node concept="lc7rE" id="7tgPrsAkt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkp" role="lcghm">
            <property role="lacIc" value="	if err := s.pre" />
          </node>
          <node concept="la8eA" id="7tgPrsAkq" role="lcghm">
            <property role="lacIc" value="{???-opName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAkr" role="lcghm">
            <property role="lacIc" value="Hook(ctx, span, &amp;event); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAks" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAku" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkx" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;400&quot;, &quot;pre-hook: &quot; + err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAky" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkA" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkD" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkH" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkG" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAkI" role="3cqZAp" />
        <!-- // ~~- DAL forwarding via Request-Reply ~~- -->
        <node concept="lc7rE" id="7tgPrsAkP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkJ" role="lcghm">
            <property role="lacIc" value="	dalSubject := s.subjectPrefix + &quot;." />
          </node>
          <node concept="la8eA" id="7tgPrsAkK" role="lcghm">
            <property role="lacIc" value="{???-nameLower}" />
          </node>
          <node concept="la8eA" id="7tgPrsAkL" role="lcghm">
            <property role="lacIc" value=".db." />
          </node>
          <node concept="la8eA" id="7tgPrsAkM" role="lcghm">
            <property role="lacIc" value="{???-opKind}" />
          </node>
          <node concept="la8eA" id="7tgPrsAkN" role="lcghm">
            <property role="lacIc" value="&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkQ" role="lcghm">
            <property role="lacIc" value="	outMsg := &amp;nats.Msg{Data: req.Data()}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkT" role="lcghm">
            <property role="lacIc" value="	outMsg.Header = core.ExtractHeaders(ctx, nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAkY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAkW" role="lcghm">
            <property role="lacIc" value="	outMsg.Header.Set(&quot;X-Business-Validated&quot;, &quot;true&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAkX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAk0" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAkZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAk3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk1" role="lcghm">
            <property role="lacIc" value="	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAk2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAk6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk4" role="lcghm">
            <property role="lacIc" value="	defer dalCancel()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAk5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAk8" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAk7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAk9" role="lcghm">
            <property role="lacIc" value="	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAla" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAle" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlc" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAld" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlf" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAli" role="lcghm">
            <property role="lacIc" value="		_ = req.RespondError(&quot;500&quot;, &quot;DAL request error: &quot; + err.Error(), nil)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAln" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAll" role="lcghm">
            <property role="lacIc" value="		return" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlo" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAls" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlt" role="lcghm">
            <property role="lacIc" value="	log.Printf(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAlu" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAlv" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAlw" role="lcghm">
            <property role="lacIc" value="{???-opKind}" />
          </node>
          <node concept="la8eA" id="7tgPrsAlx" role="lcghm">
            <property role="lacIc" value=" DAL reply: %d bytes&quot;, len(reply.Data))" />
          </node>
          <node concept="l8MVK" id="7tgPrsAly" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlB" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlA" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAlC" role="3cqZAp" />
        <!-- // ~~- Post-hook ~~- -->
        <node concept="lc7rE" id="7tgPrsAlH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlD" role="lcghm">
            <property role="lacIc" value="	responseData := s.post" />
          </node>
          <node concept="la8eA" id="7tgPrsAlE" role="lcghm">
            <property role="lacIc" value="{???-opName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAlF" role="lcghm">
            <property role="lacIc" value="Hook(ctx, span, &amp;event, reply.Data)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlI" role="lcghm">
            <property role="lacIc" value="	_ = req.Respond(responseData)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlN" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlL" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAlM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAlP" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAlO" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAlQ" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAlR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlS" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAlT" role="3cqZAp" />
        <!-- // === Pre/Post Hook Stubs === -->
        <node concept="lc7rE" id="7tgPrsAlU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlV" role="lcghm">
            <property role="lacIc" value="{???-foreach op in node.operations {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAlW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAlX" role="lcghm">
            <property role="lacIc" value="{???-string hookName = op.capitalizedName();}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAlY" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAlZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAl0" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:create) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAl9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAl1" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAl2" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAl3" role="lcghm">
            <property role="lacIc" value="Handler) pre" />
          </node>
          <node concept="la8eA" id="7tgPrsAl4" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAl5" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAl6" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAl7" role="lcghm">
            <property role="lacIc" value="CreatedEvent) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAl8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAma" role="lcghm">
            <property role="lacIc" value="	return nil" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmd" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAme" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmh" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmi" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmj" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmk" role="lcghm">
            <property role="lacIc" value="Handler) post" />
          </node>
          <node concept="la8eA" id="7tgPrsAml" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmm" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmn" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmo" role="lcghm">
            <property role="lacIc" value="CreatedEvent, data []byte) []byte {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmp" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmr" role="lcghm">
            <property role="lacIc" value="	return data" />
          </node>
          <node concept="l8MVK" id="7tgPrsAms" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmu" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmv" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmy" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmx" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmA" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAmB" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAmC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmD" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:update) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAmM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmE" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmF" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmG" role="lcghm">
            <property role="lacIc" value="Handler) pre" />
          </node>
          <node concept="la8eA" id="7tgPrsAmH" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmI" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmJ" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmK" role="lcghm">
            <property role="lacIc" value="UpdatedEvent) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmN" role="lcghm">
            <property role="lacIc" value="	return nil" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmO" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmQ" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAmR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAmU" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAmT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAm3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAmV" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAmW" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmX" role="lcghm">
            <property role="lacIc" value="Handler) post" />
          </node>
          <node concept="la8eA" id="7tgPrsAmY" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAmZ" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAm0" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAm1" role="lcghm">
            <property role="lacIc" value="UpdatedEvent, data []byte) []byte {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAm2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAm6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm4" role="lcghm">
            <property role="lacIc" value="	return data" />
          </node>
          <node concept="l8MVK" id="7tgPrsAm5" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAm9" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAm7" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAm8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnb" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAna" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnd" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAne" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAnf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAng" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:delete) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAnp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnh" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAni" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnj" role="lcghm">
            <property role="lacIc" value="Handler) pre" />
          </node>
          <node concept="la8eA" id="7tgPrsAnk" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnl" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAnm" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnn" role="lcghm">
            <property role="lacIc" value="DeletedEvent) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAno" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAns" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnq" role="lcghm">
            <property role="lacIc" value="	return nil" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnt" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnx" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAnw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAny" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAnz" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnA" role="lcghm">
            <property role="lacIc" value="Handler) post" />
          </node>
          <node concept="la8eA" id="7tgPrsAnB" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnC" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAnD" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnE" role="lcghm">
            <property role="lacIc" value="DeletedEvent, data []byte) []byte {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnF" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnH" role="lcghm">
            <property role="lacIc" value="	return data" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnI" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnK" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAnL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAnN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAnP" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnQ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAnR" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAnS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnT" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:get) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAn2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAnU" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAnV" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnW" role="lcghm">
            <property role="lacIc" value="Handler) pre" />
          </node>
          <node concept="la8eA" id="7tgPrsAnX" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAnY" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAnZ" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAn0" role="lcghm">
            <property role="lacIc" value="GetRequest) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAn5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn3" role="lcghm">
            <property role="lacIc" value="	return nil" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAn8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAn6" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAn7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoa" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAn9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAob" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAoc" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAod" role="lcghm">
            <property role="lacIc" value="Handler) post" />
          </node>
          <node concept="la8eA" id="7tgPrsAoe" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAof" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAog" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAoh" role="lcghm">
            <property role="lacIc" value="GetRequest, data []byte) []byte {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoi" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAom" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAok" role="lcghm">
            <property role="lacIc" value="	return data" />
          </node>
          <node concept="l8MVK" id="7tgPrsAol" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAop" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAon" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAor" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAoq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAos" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAot" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAou" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAov" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAow" role="lcghm">
            <property role="lacIc" value="{???-if (op.entityOperation == EntityOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAoF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAox" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAoy" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAoz" role="lcghm">
            <property role="lacIc" value="Handler) pre" />
          </node>
          <node concept="la8eA" id="7tgPrsAoA" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAoB" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAoC" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAoD" role="lcghm">
            <property role="lacIc" value="ListRequest) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoG" role="lcghm">
            <property role="lacIc" value="	return nil" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoJ" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoN" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAoM" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoO" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAoP" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAoQ" role="lcghm">
            <property role="lacIc" value="Handler) post" />
          </node>
          <node concept="la8eA" id="7tgPrsAoR" role="lcghm">
            <property role="lacIc" value="{???-hookName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAoS" role="lcghm">
            <property role="lacIc" value="Hook(ctx context.Context, span tracer.Span, event *" />
          </node>
          <node concept="la8eA" id="7tgPrsAoT" role="lcghm">
            <property role="lacIc" value="{???-name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAoU" role="lcghm">
            <property role="lacIc" value="ListRequest, data []byte) []byte {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAoZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAoX" role="lcghm">
            <property role="lacIc" value="	return data" />
          </node>
          <node concept="l8MVK" id="7tgPrsAoY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAo2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo0" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAo1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAo4" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAo3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAo5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo6" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAo7" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAo8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAo9" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsApa" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsApb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApc" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsApd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsApe" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
      </node>
    </node>
  </node>
</model>
