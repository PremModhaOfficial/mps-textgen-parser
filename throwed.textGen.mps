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
  <node concept="WtQ9Q" id="7tgPrsAhC">
    <ref role="WuzLi" to="s1:TODO-CONCEPT-ID" resolve="Relation" />
    <node concept="29tfMY" id="7tgPrsAhF" role="29tGrW">
      <node concept="3clFbS" id="7tgPrsAhG" role="2VODD2">
        <node concept="3cpWs6" id="7tgPrsAhH" role="3cqZAp">
          <node concept="Xl_RD" id="7tgPrsAhI" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="7tgPrsAhD" role="11c4hB">
      <node concept="3clFbS" id="7tgPrsAhE" role="2VODD2">
        <node concept="3clFbH" id="7tgPrsAa5" role="3cqZAp" />
        <!-- // === Package + Imports === -->
        <node concept="lc7rE" id="7tgPrsAa8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa6" role="lcghm">
            <property role="lacIc" value="package main" />
          </node>
          <node concept="l8MVK" id="7tgPrsAa7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAba" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAa9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbb" role="lcghm">
            <property role="lacIc" value="import (" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbe" role="lcghm">
            <property role="lacIc" value="	&quot;context&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbh" role="lcghm">
            <property role="lacIc" value="	&quot;encoding/json&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbi" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbk" role="lcghm">
            <property role="lacIc" value="	&quot;fmt&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbn" role="lcghm">
            <property role="lacIc" value="	&quot;time&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbr" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbs" role="lcghm">
            <property role="lacIc" value="	&quot;dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbt" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbv" role="lcghm">
            <property role="lacIc" value="	&quot;dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAby" role="lcghm">
            <property role="lacIc" value="	&quot;dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbB" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbF" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAbE" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAbG" role="3cqZAp" />
        <!-- // === Event / Request Types === -->
        <node concept="lc7rE" id="7tgPrsAbH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbI" role="lcghm">
            <property role="lacIc" value="{???-foreach op in node.operations {}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAbJ" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAbK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbL" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:assign) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAbR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbM" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAbN" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbO" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbP" role="lcghm">
            <property role="lacIc" value="AssignedEvent struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbS" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAbT" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbU" role="lcghm">
            <property role="lacIc" value="ID string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAbV" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAbW" role="lcghm">
            <property role="lacIc" value="_id&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAb5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbZ" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAb0" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb1" role="lcghm">
            <property role="lacIc" value="ID string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAb2" role="lcghm">
            <property role="lacIc" value="{???-node.to.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb3" role="lcghm">
            <property role="lacIc" value="_id&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAb8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb6" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb9" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAca" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcd" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAce" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcf" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAcg" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAch" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAci" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:remove) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAco" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcj" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAck" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcl" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcm" role="lcghm">
            <property role="lacIc" value="RemovedEvent struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcn" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcp" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcq" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcr" role="lcghm">
            <property role="lacIc" value="ID string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcs" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAct" role="lcghm">
            <property role="lacIc" value="_id&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcw" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcx" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcy" role="lcghm">
            <property role="lacIc" value="ID string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcz" role="lcghm">
            <property role="lacIc" value="{???-node.to.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcA" role="lcghm">
            <property role="lacIc" value="_id&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcD" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcG" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcK" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAcJ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcM" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAcN" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAcO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcP" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcQ" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAcR" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcS" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcT" role="lcghm">
            <property role="lacIc" value="ListRequest struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcW" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcX" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcY" role="lcghm">
            <property role="lacIc" value="ID string `json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcZ" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAc0" role="lcghm">
            <property role="lacIc" value="_id&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc3" role="lcghm">
            <property role="lacIc" value="	Limit     int       `json:&quot;limit&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAc8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc6" role="lcghm">
            <property role="lacIc" value="	Offset    int       `json:&quot;offset&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc7" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc9" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAda" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAde" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdc" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdg" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdf" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdi" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdj" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdl" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdm" role="3cqZAp" />
        <!-- // === Handler Struct + Constructor === -->
        <node concept="lc7rE" id="7tgPrsAds" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdn" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAdo" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdp" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdq" role="lcghm">
            <property role="lacIc" value="Handler struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdt" role="lcghm">
            <property role="lacIc" value="	publisher     *nats.Publisher" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdu" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdw" role="lcghm">
            <property role="lacIc" value="	subjectPrefix string" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdx" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdz" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdA" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdD" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAdC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdE" role="lcghm">
            <property role="lacIc" value="func New" />
          </node>
          <node concept="la8eA" id="7tgPrsAdF" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdG" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdH" role="lcghm">
            <property role="lacIc" value="Handler(pub *nats.Publisher, subjectPrefix string) *" />
          </node>
          <node concept="la8eA" id="7tgPrsAdI" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdJ" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdK" role="lcghm">
            <property role="lacIc" value="Handler {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdN" role="lcghm">
            <property role="lacIc" value="	return &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAdO" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdP" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdQ" role="lcghm">
            <property role="lacIc" value="Handler{" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdT" role="lcghm">
            <property role="lacIc" value="		publisher:     pub," />
          </node>
          <node concept="l8MVK" id="7tgPrsAdU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdW" role="lcghm">
            <property role="lacIc" value="		subjectPrefix: subjectPrefix," />
          </node>
          <node concept="l8MVK" id="7tgPrsAdX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdZ" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAd0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd2" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAd3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd6" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAd5" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAd7" role="3cqZAp" />
        <!-- // === Handler Methods === -->
        <node concept="lc7rE" id="7tgPrsAd8" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd9" role="lcghm">
            <property role="lacIc" value="{???-foreach op in node.operations {}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAea" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAei" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeb" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAec" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAed" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAee" role="lcghm">
            <property role="lacIc" value="Handler) Handle" />
          </node>
          <node concept="la8eA" id="7tgPrsAef" role="lcghm">
            <property role="lacIc" value="{???-op.capitalizedName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeg" role="lcghm">
            <property role="lacIc" value="(ctx context.Context, msg *core.Message) error {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAej" role="lcghm">
            <property role="lacIc" value="	ctx, span := tracer.StartConsumer(ctx, &quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAek" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAel" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAem" role="lcghm">
            <property role="lacIc" value=".Handle" />
          </node>
          <node concept="la8eA" id="7tgPrsAen" role="lcghm">
            <property role="lacIc" value="{???-op.capitalizedName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeo" role="lcghm">
            <property role="lacIc" value="&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAep" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAet" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAer" role="lcghm">
            <property role="lacIc" value="	defer span.End()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAes" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAew" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeu" role="lcghm">
            <property role="lacIc" value="	ctx = core.InjectContext(ctx, msg.Headers)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAev" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAey" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAex" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAez" role="3cqZAp" />
        <!-- // ~~- Unmarshal ~~- -->
        <node concept="lc7rE" id="7tgPrsAeA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeB" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:assign) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeC" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAeD" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeE" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeF" role="lcghm">
            <property role="lacIc" value="AssignedEvent" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeG" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeJ" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeL" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:remove) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeM" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAeN" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeO" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeP" role="lcghm">
            <property role="lacIc" value="RemovedEvent" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeT" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeV" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAe1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeW" role="lcghm">
            <property role="lacIc" value="	var event " />
          </node>
          <node concept="la8eA" id="7tgPrsAeX" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeY" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeZ" role="lcghm">
            <property role="lacIc" value="ListRequest" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAe2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe3" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAe4" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAe7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe5" role="lcghm">
            <property role="lacIc" value="	if err := json.Unmarshal(msg.Data, &amp;event); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfa" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe8" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfb" role="lcghm">
            <property role="lacIc" value="		return err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfg" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfe" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAff" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfi" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfh" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfj" role="3cqZAp" />
        <!-- // ~~- Validation ~~- -->
        <node concept="lc7rE" id="7tgPrsAfk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfl" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:assign || op.relationOperation == RelationOperation:remove) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfm" role="lcghm">
            <property role="lacIc" value="	if event." />
          </node>
          <node concept="la8eA" id="7tgPrsAfn" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfo" role="lcghm">
            <property role="lacIc" value="ID == &quot;&quot; || event." />
          </node>
          <node concept="la8eA" id="7tgPrsAfp" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfq" role="lcghm">
            <property role="lacIc" value="ID == &quot;&quot; {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfz" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAft" role="lcghm">
            <property role="lacIc" value="		err := fmt.Errorf(&quot;invalid data: missing " />
          </node>
          <node concept="la8eA" id="7tgPrsAfu" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfv" role="lcghm">
            <property role="lacIc" value=" or " />
          </node>
          <node concept="la8eA" id="7tgPrsAfw" role="lcghm">
            <property role="lacIc" value="{???-node.to.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfx" role="lcghm">
            <property role="lacIc" value=" ID&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfy" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfA" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfD" role="lcghm">
            <property role="lacIc" value="		return err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfG" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfK" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfM" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfN" role="lcghm">
            <property role="lacIc" value="	if event." />
          </node>
          <node concept="la8eA" id="7tgPrsAfO" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfP" role="lcghm">
            <property role="lacIc" value="ID == &quot;&quot; {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfS" role="lcghm">
            <property role="lacIc" value="		err := fmt.Errorf(&quot;invalid request: missing " />
          </node>
          <node concept="la8eA" id="7tgPrsAfT" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfU" role="lcghm">
            <property role="lacIc" value=" ID&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfX" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfY" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAf2" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf0" role="lcghm">
            <property role="lacIc" value="		return err" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf1" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAf5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf3" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf4" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAf6" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf7" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAf8" role="3cqZAp" />
        <!-- // ~~- Span Attributes ~~- -->
        <node concept="lc7rE" id="7tgPrsAga" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAf9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgb" role="lcghm">
            <property role="lacIc" value="	span.SetAttributes(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgc" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAge" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgf" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:assign || op.relationOperation == RelationOperation:remove) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgg" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAgh" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgi" role="lcghm">
            <property role="lacIc" value=".id&quot;, event." />
          </node>
          <node concept="la8eA" id="7tgPrsAgj" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgk" role="lcghm">
            <property role="lacIc" value="ID)," />
          </node>
          <node concept="l8MVK" id="7tgPrsAgl" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgn" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAgo" role="lcghm">
            <property role="lacIc" value="{???-node.to.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgp" role="lcghm">
            <property role="lacIc" value=".id&quot;, event." />
          </node>
          <node concept="la8eA" id="7tgPrsAgq" role="lcghm">
            <property role="lacIc" value="{???-node.to.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgr" role="lcghm">
            <property role="lacIc" value="ID)," />
          </node>
          <node concept="l8MVK" id="7tgPrsAgs" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgu" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgv" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgw" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgx" role="lcghm">
            <property role="lacIc" value="{???-if (op.relationOperation == RelationOperation:list) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgE" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgy" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAgz" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgA" role="lcghm">
            <property role="lacIc" value=".id&quot;, event." />
          </node>
          <node concept="la8eA" id="7tgPrsAgB" role="lcghm">
            <property role="lacIc" value="{???-node.from.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgC" role="lcghm">
            <property role="lacIc" value="ID)," />
          </node>
          <node concept="l8MVK" id="7tgPrsAgD" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgG" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgH" role="lcghm">
            <property role="lacIc" value="		tracer.StringAttr(&quot;tenant.id&quot;, msg.Headers.Get(core.HeaderTenantID))," />
          </node>
          <node concept="l8MVK" id="7tgPrsAgI" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgK" role="lcghm">
            <property role="lacIc" value="	)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgL" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgO" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAgN" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgP" role="3cqZAp" />
        <!-- // ~~- Publish to DAL ~~- -->
        <node concept="lc7rE" id="7tgPrsAgS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgQ" role="lcghm">
            <property role="lacIc" value="	outMsg := core.NewMessage(msg.Data)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAg1" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgT" role="lcghm">
            <property role="lacIc" value="	outMsg.Subject = s.subjectPrefix + &quot;." />
          </node>
          <node concept="la8eA" id="7tgPrsAgU" role="lcghm">
            <property role="lacIc" value="{???-node.from.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgV" role="lcghm">
            <property role="lacIc" value="." />
          </node>
          <node concept="la8eA" id="7tgPrsAgW" role="lcghm">
            <property role="lacIc" value="{???-node.to.name.toLowerCase()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgX" role="lcghm">
            <property role="lacIc" value=".db." />
          </node>
          <node concept="la8eA" id="7tgPrsAgY" role="lcghm">
            <property role="lacIc" value="{???-op.relationOperation.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgZ" role="lcghm">
            <property role="lacIc" value="&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAg0" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAg4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg2" role="lcghm">
            <property role="lacIc" value="	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAg3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAg7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAg5" role="lcghm">
            <property role="lacIc" value="	outMsg.Headers.Set(&quot;X-Business-Validated&quot;, &quot;true&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAg6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAg9" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAg8" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhc" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAha" role="lcghm">
            <property role="lacIc" value="	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhb" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhd" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhe" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhi" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhg" role="lcghm">
            <property role="lacIc" value="		return fmt.Errorf(&quot;publish error: %w&quot;, err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhh" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhj" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAho" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhm" role="lcghm">
            <property role="lacIc" value="	return nil" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhn" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAhr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhp" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAhq" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAht" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAhs" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAhu" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAhv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhw" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAhx" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAhy" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhz" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAhA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAhB" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
      </node>
    </node>
  </node>
</model>
