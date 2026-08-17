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
  <node concept="WtQ9Q" id="7tgPrsAgJ">
    <ref role="WuzLi" to="s1:TODO-CONCEPT-ID" resolve="NatsService" />
    <node concept="29tfMY" id="7tgPrsAgM" role="29tGrW">
      <node concept="3clFbS" id="7tgPrsAgN" role="2VODD2">
        <node concept="3cpWs6" id="7tgPrsAgO" role="3cqZAp">
          <node concept="Xl_RD" id="7tgPrsAgP" role="3cqZAk">
            <property role="Xl_RC" value="generated" />
          </node>
        </node>
      </node>
    </node>
    <node concept="11bSqf" id="7tgPrsAgK" role="11c4hB">
      <node concept="3clFbS" id="7tgPrsAgL" role="2VODD2">
        <node concept="3clFbH" id="7tgPrsAa5" role="3cqZAp" />
        <!-- // === 0. Package and Imports === -->
        <node concept="lc7rE" id="7tgPrsAba" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAa6" role="lcghm">
            <property role="lacIc" value="package " />
          </node>
          <node concept="la8eA" id="7tgPrsAa7" role="lcghm">
            <property role="lacIc" value="{???-node.packageName}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAa8" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAa9" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAbb" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAbe" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbc" role="lcghm">
            <property role="lacIc" value="import (" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbf" role="lcghm">
            <property role="lacIc" value="	&quot;context&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbi" role="lcghm">
            <property role="lacIc" value="	&quot;encoding/json&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbl" role="lcghm">
            <property role="lacIc" value="	&quot;fmt&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbo" role="lcghm">
            <property role="lacIc" value="	&quot;time&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbp" role="lcghm" />
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
            <property role="lacIc" value="	&quot;go.opentelemetry.io/otel/attribute&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbB" role="lcghm">
            <property role="lacIc" value="	&quot;go.opentelemetry.io/otel/trace&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbE" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbF" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbG" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAbI" role="3cqZAp" />
        <!-- // === 1. Data Access Layer === -->
        <node concept="lc7rE" id="7tgPrsAbL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbJ" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbM" role="lcghm">
            <property role="lacIc" value="// 1. Data Access Layer" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbS" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbP" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbQ" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAbR" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAbV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbT" role="lcghm">
            <property role="lacIc" value="type DataAccessLayer interface {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAbU" role="lcghm" />
        </node>
        <!-- // Dynamic DAL generation based on models -->
        <node concept="lc7rE" id="7tgPrsAbW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbX" role="lcghm">
            <property role="lacIc" value="{???-foreach model in node.models {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAb4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAbY" role="lcghm">
            <property role="lacIc" value="	Get" />
          </node>
          <node concept="la8eA" id="7tgPrsAbZ" role="lcghm">
            <property role="lacIc" value="{???-model.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb0" role="lcghm">
            <property role="lacIc" value="ByID(ctx context.Context, id string) (*" />
          </node>
          <node concept="la8eA" id="7tgPrsAb1" role="lcghm">
            <property role="lacIc" value="{???-model.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAb2" role="lcghm">
            <property role="lacIc" value=", error)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAb5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb6" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAca" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAb7" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAb8" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAb9" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAcb" role="3cqZAp" />
        <!-- // === 2. Domain Models === -->
        <node concept="lc7rE" id="7tgPrsAce" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcc" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAch" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcf" role="lcghm">
            <property role="lacIc" value="// 2. Domain Models (SQL Structures)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcl" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAci" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcj" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAck" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcn" role="lcghm">
            <property role="lacIc" value="{???-foreach model in node.models {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcs" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAco" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAcp" role="lcghm">
            <property role="lacIc" value="{???-model.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcq" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcr" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAct" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcu" role="lcghm">
            <property role="lacIc" value="{???-foreach field in model.fields {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcv" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcw" role="lcghm">
            <property role="lacIc" value="{???-field.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcx" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAcy" role="lcghm">
            <property role="lacIc" value="{???-field.goType}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcz" role="lcghm">
            <property role="lacIc" value="	`json:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcA" role="lcghm">
            <property role="lacIc" value="{???-field.jsonName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcB" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAcC" role="lcghm">
            <property role="lacIc" value="{???-field.dbName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAcD" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcE" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcG" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcH" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAcL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcI" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcJ" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAcK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcM" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcN" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAcO" role="3cqZAp" />
        <!-- // === 3. Event Payloads === -->
        <node concept="lc7rE" id="7tgPrsAcR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcP" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcS" role="lcghm">
            <property role="lacIc" value="// 3. Event Payloads (NATS Messages)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcY" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAcV" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAcW" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAcX" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAcZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc0" role="lcghm">
            <property role="lacIc" value="{???-foreach event in node.events {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAc5" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc1" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAc2" role="lcghm">
            <property role="lacIc" value="{???-event.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAc3" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc4" role="lcghm" />
        </node>
        <!-- // Use modern collection iteration to reduce duplication for event fields -->
        <node concept="lc7rE" id="7tgPrsAc7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc6" role="lcghm">
            <property role="lacIc" value="{???-event.fields.forEach(~it =&gt; it.name + &quot; &quot; + it.goType + &quot; `json:\&quot;&quot; + it.jsonName + &quot;\&quot;`\n&quot;)}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAda" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAc8" role="lcghm">
            <property role="lacIc" value="	Timestamp time.Time `json:&quot;timestamp&quot;`" />
          </node>
          <node concept="l8MVK" id="7tgPrsAc9" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAde" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdb" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdc" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAdd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdf" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdg" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdh" role="3cqZAp" />
        <!-- // === 4. Service Implementation === -->
        <node concept="lc7rE" id="7tgPrsAdk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdi" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdn" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdl" role="lcghm">
            <property role="lacIc" value="// 4. Service Implementation" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdm" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdr" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdo" role="lcghm">
            <property role="lacIc" value="// ==========================================" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdp" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAdq" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAds" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAdt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdu" role="lcghm">
            <property role="lacIc" value="{???-string serviceInterface = node.name + &quot;Service&quot;;}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdv" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdw" role="lcghm">
            <property role="lacIc" value="{???-string serviceStruct = node.name + &quot;Service&quot;;}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAdx" role="3cqZAp" />
        <!-- // Interface -->
        <node concept="lc7rE" id="7tgPrsAdC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdy" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAdz" role="lcghm">
            <property role="lacIc" value="{???-node.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdA" role="lcghm">
            <property role="lacIc" value="Service interface {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdB" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdE" role="lcghm">
            <property role="lacIc" value="{???-foreach method in node.methods {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdI" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdF" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
          <node concept="la8eA" id="7tgPrsAdG" role="lcghm">
            <property role="lacIc" value="{???-method.signature()}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdH" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdJ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdK" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAdO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdL" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdM" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAdN" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAdP" role="3cqZAp" />
        <!-- // Struct -->
        <node concept="lc7rE" id="7tgPrsAdU" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdQ" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
          <node concept="la8eA" id="7tgPrsAdR" role="lcghm">
            <property role="lacIc" value="{???-node.camelName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAdS" role="lcghm">
            <property role="lacIc" value="Service struct {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdT" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAdX" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdV" role="lcghm">
            <property role="lacIc" value="	publisher *nats.Publisher" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdW" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAdY" role="lcghm">
            <property role="lacIc" value="	tracer    trace.Tracer" />
          </node>
          <node concept="l8MVK" id="7tgPrsAdZ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd1" role="lcghm">
            <property role="lacIc" value="	dal       DataAccessLayer" />
          </node>
          <node concept="l8MVK" id="7tgPrsAd2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAd7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd4" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAd5" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAd6" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAd8" role="3cqZAp" />
        <!-- // Constructor -->
        <node concept="lc7rE" id="7tgPrsAef" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAd9" role="lcghm">
            <property role="lacIc" value="func New" />
          </node>
          <node concept="la8eA" id="7tgPrsAea" role="lcghm">
            <property role="lacIc" value="{???-node.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeb" role="lcghm">
            <property role="lacIc" value="Service(pub *nats.Publisher, tracer trace.Tracer, dal DataAccessLayer) " />
          </node>
          <node concept="la8eA" id="7tgPrsAec" role="lcghm">
            <property role="lacIc" value="{???-node.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAed" role="lcghm">
            <property role="lacIc" value="Service {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAee" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAek" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeg" role="lcghm">
            <property role="lacIc" value="	return &amp;" />
          </node>
          <node concept="la8eA" id="7tgPrsAeh" role="lcghm">
            <property role="lacIc" value="{???-node.camelName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAei" role="lcghm">
            <property role="lacIc" value="Service{" />
          </node>
          <node concept="l8MVK" id="7tgPrsAej" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAen" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAel" role="lcghm">
            <property role="lacIc" value="		publisher: pub," />
          </node>
          <node concept="l8MVK" id="7tgPrsAem" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeq" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeo" role="lcghm">
            <property role="lacIc" value="		tracer:    tracer," />
          </node>
          <node concept="l8MVK" id="7tgPrsAep" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAet" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAer" role="lcghm">
            <property role="lacIc" value="		dal:       dal," />
          </node>
          <node concept="l8MVK" id="7tgPrsAes" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAew" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeu" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAev" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAex" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAey" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAez" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAeB" role="3cqZAp" />
        <!-- // Methods -->
        <node concept="lc7rE" id="7tgPrsAeC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeD" role="lcghm">
            <property role="lacIc" value="{???-foreach method in node.methods {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAeK" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeE" role="lcghm">
            <property role="lacIc" value="func (s *" />
          </node>
          <node concept="la8eA" id="7tgPrsAeF" role="lcghm">
            <property role="lacIc" value="{???-node.camelName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeG" role="lcghm">
            <property role="lacIc" value="Service) " />
          </node>
          <node concept="la8eA" id="7tgPrsAeH" role="lcghm">
            <property role="lacIc" value="{???-method.signatureWithNames()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeI" role="lcghm">
            <property role="lacIc" value=" {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeJ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAeL" role="3cqZAp" />
        <!-- // 4a. OTEL Span setup -->
        <node concept="lc7rE" id="7tgPrsAeO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeM" role="lcghm">
            <property role="lacIc" value="	// Start OTEL Span" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeV" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeP" role="lcghm">
            <property role="lacIc" value="	ctx, span := s.tracer.Start(ctx, &quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAeQ" role="lcghm">
            <property role="lacIc" value="{???-node.pascalName()}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeR" role="lcghm">
            <property role="lacIc" value="Service." />
          </node>
          <node concept="la8eA" id="7tgPrsAeS" role="lcghm">
            <property role="lacIc" value="{???-method.name}" />
          </node>
          <node concept="la8eA" id="7tgPrsAeT" role="lcghm">
            <property role="lacIc" value="&quot;)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeU" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAeZ" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAeW" role="lcghm">
            <property role="lacIc" value="	defer span.End()" />
          </node>
          <node concept="l8MVK" id="7tgPrsAeX" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAeY" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAe0" role="3cqZAp" />
        <!-- // Dynamic span attributes using modern list joining -->
        <node concept="lc7rE" id="7tgPrsAe3" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe1" role="lcghm">
            <property role="lacIc" value="	span.SetAttributes(" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe2" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAe7" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe4" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
          <node concept="la8eA" id="7tgPrsAe5" role="lcghm">
            <property role="lacIc" value="{???-method.spanAttributes.select(~it =&gt; &quot;attribute.String(\&quot;&quot; + it.key + &quot;\&quot;, &quot; + it.variable + &quot;)&quot;).join(&quot;,\n\t\t&quot;)}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe6" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfb" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAe8" role="lcghm">
            <property role="lacIc" value="	)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAe9" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAfa" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfc" role="3cqZAp" />
        <!-- // 4b. DAL Validations (Dynamic based on method.dalChecks) -->
        <node concept="lc7rE" id="7tgPrsAfd" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfe" role="lcghm">
            <property role="lacIc" value="{???-if (method.hasDalChecks) {}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAff" role="lcghm">
            <property role="lacIc" value="	// Validate against DAL" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfj" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfi" role="lcghm">
            <property role="lacIc" value="{???-method.dalChecks.forEach(~it =&gt; &quot;_, err := s.dal.Get&quot; + it.model + &quot;ByID(ctx, &quot; + it.varName + &quot;)\n\tif err != nil {\n\t\tspan.RecordError(err)\n\t\treturn fmt.Errorf(\&quot;&quot; + it.model + &quot; not found: %w\&quot;, err)\n\t}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfl" role="3cqZAp">
          <node concept="l8MVK" id="7tgPrsAfk" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfm" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfn" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAfo" role="3cqZAp" />
        <!-- // 4c. Event Creation -->
        <node concept="lc7rE" id="7tgPrsAft" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfp" role="lcghm">
            <property role="lacIc" value="	event := " />
          </node>
          <node concept="la8eA" id="7tgPrsAfq" role="lcghm">
            <property role="lacIc" value="{???-method.eventName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfr" role="lcghm">
            <property role="lacIc" value="{" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfs" role="lcghm" />
        </node>
        <!-- // Use modern syntax to map method params to event fields -->
        <node concept="lc7rE" id="7tgPrsAfx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfu" role="lcghm">
            <property role="lacIc" value=" method.eventMappings with &quot;," />
          </node>
          <node concept="l8MVK" id="7tgPrsAfv" role="lcghm" />
          <node concept="la8eA" id="7tgPrsAfw" role="lcghm">
            <property role="lacIc" value="		&quot; " />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAfA" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfy" role="lcghm">
            <property role="lacIc" value="," />
          </node>
          <node concept="l8MVK" id="7tgPrsAfz" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfD" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfB" role="lcghm">
            <property role="lacIc" value="		Timestamp: time.Now().UTC()," />
          </node>
          <node concept="l8MVK" id="7tgPrsAfC" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfE" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfF" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAfG" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAfI" role="3cqZAp" />
        <!-- // 4d. Publish logic -->
        <node concept="lc7rE" id="7tgPrsAfL" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfJ" role="lcghm">
            <property role="lacIc" value="	payload, err := json.Marshal(event)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfK" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfO" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfM" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfN" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfR" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfP" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfQ" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAfW" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfS" role="lcghm">
            <property role="lacIc" value="		return fmt.Errorf(&quot;failed to marshal " />
          </node>
          <node concept="la8eA" id="7tgPrsAfT" role="lcghm">
            <property role="lacIc" value="{???-method.eventName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAfU" role="lcghm">
            <property role="lacIc" value=": %w&quot;, err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfV" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAf0" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAfX" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAfY" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAfZ" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAf1" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAf4" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf2" role="lcghm">
            <property role="lacIc" value="	msg := core.NewMessage(payload)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf3" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAga" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAf5" role="lcghm">
            <property role="lacIc" value="	msg.Subject = &quot;" />
          </node>
          <node concept="la8eA" id="7tgPrsAf6" role="lcghm">
            <property role="lacIc" value="{???-method.subjectName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAf7" role="lcghm">
            <property role="lacIc" value="&quot;" />
          </node>
          <node concept="l8MVK" id="7tgPrsAf8" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAf9" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgb" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAge" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgc" role="lcghm">
            <property role="lacIc" value="	// Publish to NATS" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgd" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgh" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgf" role="lcghm">
            <property role="lacIc" value="	if err := s.publisher.Publish(ctx, msg.Subject, msg); err != nil {" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgg" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgk" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgi" role="lcghm">
            <property role="lacIc" value="		span.RecordError(err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgj" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgp" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgl" role="lcghm">
            <property role="lacIc" value="		return fmt.Errorf(&quot;failed to publish " />
          </node>
          <node concept="la8eA" id="7tgPrsAgm" role="lcghm">
            <property role="lacIc" value="{???-method.eventName}" />
          </node>
          <node concept="la8eA" id="7tgPrsAgn" role="lcghm">
            <property role="lacIc" value=" event: %w&quot;, err)" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgo" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgt" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgq" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgr" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAgs" role="lcghm" />
        </node>
        <node concept="3clFbH" id="7tgPrsAgu" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAgx" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgv" role="lcghm">
            <property role="lacIc" value="	return nil" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgw" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgB" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgy" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
          <node concept="l8MVK" id="7tgPrsAgz" role="lcghm" />
          <node concept="l8MVK" id="7tgPrsAgA" role="lcghm" />
        </node>
        <node concept="lc7rE" id="7tgPrsAgC" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgD" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="3clFbH" id="7tgPrsAgE" role="3cqZAp" />
        <node concept="lc7rE" id="7tgPrsAgF" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgG" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
        <node concept="lc7rE" id="7tgPrsAgH" role="3cqZAp">
          <node concept="la8eA" id="7tgPrsAgI" role="lcghm">
            <property role="lacIc" value="{???-}}" />
          </node>
        </node>
      </node>
    </node>
  </node>
</model>
