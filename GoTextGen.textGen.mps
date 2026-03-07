<?xml version="1.0" encoding="UTF-8"?>
<model ref="r:TODO-MODEL-UUID(YourLang.textGen)">
  <persistence version="9" />
  <languages>
    <use id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen" version="1" />
    <devkit ref="fa73d85a-ac7f-447b-846c-fcdc41caa600(jetbrains.mps.devkit.aspect.textgen)" />
  </languages>
  <imports>
    <import index="s1" ref="r:TODO-STRUCTURE-UUID(YourLang.structure)" implicit="true" />
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
  <node concept="WtQ9Q" id="#Itg_2161">
    <ref role="WuzLi" to="s1:TODO-CONCEPT-ID" resolve="Code" />
    <node concept="11bSqf" id="#Igt_2162" role="11c4hB">
      <node concept="3clFbS" id="#Isl_2163" role="2VODD2">
        <!-- VAR: node&lt;Code&gt; n = node; -->
        <!-- VAR: node&lt;Models&gt; model = n.model; -->
        <!-- VAR: node&lt;Infra&gt; infra = n.infra; -->
        <node concept="3clFbH" id="#Is_0001" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_0002" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0003" role="lcghm">
            <property role="lacIc" value="package main" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0004" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0005" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0006" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0007" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0008" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0009" role="lcghm">
            <property role="lacIc" value="import (" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0010" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0011" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0012" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0013" role="lcghm">
            <property role="lacIc" value="	&quot;database/sql&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0014" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0015" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0016" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0017" role="lcghm">
            <property role="lacIc" value="	_ &quot;embed&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0018" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0019" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0020" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0021" role="lcghm">
            <property role="lacIc" value="	&quot;encoding/json&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0022" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0023" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0024" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0025" role="lcghm">
            <property role="lacIc" value="	&quot;fmt&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0026" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0027" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0028" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0029" role="lcghm">
            <property role="lacIc" value="	&quot;log&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0030" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0031" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0032" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0033" role="lcghm">
            <property role="lacIc" value="	&quot;net/http&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0034" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0035" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0036" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0037" role="lcghm">
            <property role="lacIc" value="	&quot;os&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0038" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0039" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0040" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0041" role="lcghm">
            <property role="lacIc" value="	&quot;strconv&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0042" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0043" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0044" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0045" role="lcghm">
            <property role="lacIc" value="	&quot;time&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0046" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0047" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0048" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0049" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0050" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0051" role="lcghm">
            <property role="lacIc" value="	_ &quot;github.com/lib/pq&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0052" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0053" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0054" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0055" role="lcghm">
            <property role="lacIc" value="	httpSwagger &quot;github.com/swaggo/http-swagger&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0056" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0057" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0058" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0059" role="lcghm">
            <property role="lacIc" value="	_ &quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0060" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0061" role="lcghm">
            <property role="lacIc" value="▶infra.modulePath◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0062" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0063" role="lcghm">
            <property role="lacIc" value="/docs&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0064" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0065" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0066" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0067" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0068" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0069" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0070" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0071" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0072" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0073" role="lcghm">
            <property role="lacIc" value="//go:embed user_management_init.sql" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0074" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0075" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0076" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0077" role="lcghm">
            <property role="lacIc" value="var migrationSQL string" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0078" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0079" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0080" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0081" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0082" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0083" role="lcghm">
            <property role="lacIc" value="// @title         " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0084" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0085" role="lcghm">
            <property role="lacIc" value="▶model.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0086" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0087" role="lcghm">
            <property role="lacIc" value=" API" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0088" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0089" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0090" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0091" role="lcghm">
            <property role="lacIc" value="// @version       1.0" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0092" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0093" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0094" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0095" role="lcghm">
            <property role="lacIc" value="// @description   CRUD service for " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0096" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0097" role="lcghm">
            <property role="lacIc" value="▶model.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0098" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0099" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0100" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0101" role="lcghm">
            <property role="lacIc" value="// @host          localhost:" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0102" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0103" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0104" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0105" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0106" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0107" role="lcghm">
            <property role="lacIc" value="// @BasePath      /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0108" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0109" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0110" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0111" role="lcghm">
            <property role="lacIc" value="// @schemes       http" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0112" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0113" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0114" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0115" role="lcghm">
            <property role="lacIc" value="// @produce       json" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0116" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0117" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0118" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0119" role="lcghm">
            <property role="lacIc" value="// @consumes      json" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0120" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0121" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0122" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0123" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0124" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0125" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0126" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0127" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0128" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0129" role="lcghm">
            <property role="lacIc" value="// Models" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0130" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0131" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0132" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0133" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0134" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0135" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0136" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0137" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0138" role="3cqZAp" />
        <!-- // ~~~~ Regular schema structs ~~~~ -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <node concept="3clFbH" id="#Is_0139" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_0140" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0141" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0142" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0143" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0144" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0145" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0146" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0147" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0148" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0149" role="lcghm">
            <property role="lacIc" value="	ID int64 `json:&quot;id&quot; db:&quot;id&quot; example:&quot;1&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0150" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0151" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0152" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0153" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0154" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0155" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0156" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0157" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0158" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0159" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0160" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0161" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0162" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0163" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0164" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0165" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0166" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0167" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0168" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0169" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0170" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0171" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0172" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0173" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0174" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0175" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0176" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0177" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0178" role="3cqZAp" />
        <!-- // MarshalJSON if Sensitive -->
        <!-- VAR: boolean hasSensitive = false; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.Sensitive) { hasSensitive = true; } -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_0179" role="3cqZAp" />
        <!-- // found the sensitive parst -->
        <!-- ─── IF: if (hasSensitive) ─── -->
        <node concept="lc7rE" id="#Ia_0180" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0181" role="lcghm">
            <property role="lacIc" value="func (u " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0182" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0183" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0184" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0185" role="lcghm">
            <property role="lacIc" value=") MarshalJSON() ([]byte, error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0186" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0187" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0188" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0189" role="lcghm">
            <property role="lacIc" value="	type Alias " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0190" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0191" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0192" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0193" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0194" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0195" role="lcghm">
            <property role="lacIc" value="	return json.Marshal(&amp;struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0196" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0197" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0198" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0199" role="lcghm">
            <property role="lacIc" value="		Alias" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0200" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0201" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.Sensitive) ─── -->
        <node concept="lc7rE" id="#Ia_0202" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0203" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0204" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0205" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0206" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0207" role="lcghm">
            <property role="lacIc" value=" string `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0208" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0209" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0210" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0211" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0212" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0213" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0214" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0215" role="lcghm">
            <property role="lacIc" value="	}{" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0216" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0217" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0218" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0219" role="lcghm">
            <property role="lacIc" value="		Alias: (Alias)(u)," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0220" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0221" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.Sensitive) ─── -->
        <node concept="lc7rE" id="#Ia_0222" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0223" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0224" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0225" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0226" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0227" role="lcghm">
            <property role="lacIc" value=": &quot;[REDACTED]&quot;," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0228" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0229" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0230" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0231" role="lcghm">
            <property role="lacIc" value="	})" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0232" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0233" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0234" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0235" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0236" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0237" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0238" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0239" role="lcghm" />
        </node>
        <!-- END IF -->
        <node concept="3clFbH" id="#Is_0240" role="3cqZAp" />
        <node concept="3clFbH" id="#Is_0241" role="3cqZAp" />
        <!-- // Create struct -->
        <node concept="lc7rE" id="#Ia_0242" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0243" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0244" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0245" role="lcghm">
            <property role="lacIc" value="▶schema.createStructName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0246" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0247" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0248" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0249" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <node concept="lc7rE" id="#Ia_0250" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0251" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0252" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0253" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0254" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0255" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0256" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0257" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0258" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0259" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0260" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0261" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0262" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0263" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0264" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0265" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0266" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0267" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0268" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0269" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0270" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0271" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_0272" role="3cqZAp" />
        <!-- // ~~~~ Join schema structs ~~~~ -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <node concept="lc7rE" id="#Ia_0273" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0274" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0275" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0276" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0277" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0278" role="lcghm">
            <property role="lacIc" value=" struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0279" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0280" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="#Ia_0281" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0282" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0283" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0284" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0285" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0286" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0287" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0288" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0289" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0290" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0291" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0292" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0293" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0294" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0295" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0296" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0297" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0298" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0299" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0300" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0301" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0302" role="lcghm">
            <property role="lacIc" value=" " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0303" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0304" role="lcghm">
            <property role="lacIc" value="▶f.goType()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0305" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0306" role="lcghm">
            <property role="lacIc" value=" `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0307" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0308" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0309" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0310" role="lcghm">
            <property role="lacIc" value="&quot; db:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0311" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0312" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0313" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0314" role="lcghm">
            <property role="lacIc" value="&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0315" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0316" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0317" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0318" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0319" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0320" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0321" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0322" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0323" role="3cqZAp" />
        <!-- VAR: int refCount = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- ─── IF: if (refCount == 1) ─── -->
        <node concept="lc7rE" id="#Ia_0324" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0325" role="lcghm">
            <property role="lacIc" value="type Assign" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0326" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0327" role="lcghm">
            <property role="lacIc" value="▶fr.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0328" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0329" role="lcghm">
            <property role="lacIc" value="Body struct {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0330" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0331" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0332" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0333" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0334" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0335" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0336" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0337" role="lcghm">
            <property role="lacIc" value=" int64 `json:&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0338" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0339" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0340" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0341" role="lcghm">
            <property role="lacIc" value="&quot; binding:&quot;required&quot;`" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0342" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0343" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0344" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0345" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0346" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0347" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0348" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0349" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- ASSIGN: refCount = refCount + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_0350" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — regular -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="#Is_0351" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_0352" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0353" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0354" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0355" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0356" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0357" role="lcghm">
            <property role="lacIc" value="// Repositories" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0358" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0359" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0360" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0361" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0362" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0363" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0364" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0365" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0366" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="3clFbH" id="#Is_0367" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_0368" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0369" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0370" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0371" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0372" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0373" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0374" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0375" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0376" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0377" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0378" role="3cqZAp" />
        <!-- // Create -->
        <node concept="lc7rE" id="#Ia_0379" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0380" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0381" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0382" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0383" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0384" role="lcghm">
            <property role="lacIc" value=") Create(u *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0385" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0386" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0387" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0388" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0389" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0390" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0391" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0392" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0393" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0394" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0395" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0396" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0397" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0398" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0399" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0400" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <!-- VAR: int idx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <!-- IF: if (idx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0401" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0402" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="#Ia_0403" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0404" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- ASSIGN: idx = idx + 1; -->
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0405" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0406" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0407" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0408" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0409" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0410" role="lcghm">
            <property role="lacIc" value="		 VALUES (" />
          </node>
        </node>
        <!-- ═══ FOR: for (int i = 1; i &lt;= idx; i++) ═══ -->
        <!-- IF: if (i &gt; 1) -->
        <node concept="lc7rE" id="#Ia_0411" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0412" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="#Ia_0413" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0414" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0415" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0416" role="lcghm">
            <property role="lacIc" value="▶i◀" />
          </node>
        </node>
        <!-- END FOR -->
        <node concept="lc7rE" id="#Ia_0417" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0418" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0419" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0420" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0421" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0422" role="lcghm">
            <property role="lacIc" value="		 RETURNING id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType == DataType.timestamp) -->
        <node concept="lc7rE" id="#Ia_0423" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0424" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0425" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0426" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0427" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0428" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0429" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0430" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0431" role="3cqZAp" />
        <!-- // non timestapm -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType != DataType.timestamp) -->
        <node concept="lc7rE" id="#Ia_0432" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0433" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0434" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0435" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0436" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0437" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0438" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0439" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_0440" role="3cqZAp" />
        <!-- // scanning -->
        <node concept="lc7rE" id="#Ia_0441" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0442" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType == DataType.timestamp) -->
        <node concept="lc7rE" id="#Ia_0443" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0444" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0445" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0446" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0447" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0448" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0449" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0450" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0451" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0452" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0453" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0454" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0455" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0456" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0457" role="3cqZAp" />
        <!-- // GetByID -->
        <node concept="lc7rE" id="#Ia_0458" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0459" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0460" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0461" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0462" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0463" role="lcghm">
            <property role="lacIc" value=") GetByID(id int64) (*" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0464" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0465" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0466" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0467" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0468" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0469" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0470" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0471" role="lcghm">
            <property role="lacIc" value="	u := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0472" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0473" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0474" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0475" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0476" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0477" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0478" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0479" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0480" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0481" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0482" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0483" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0484" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0485" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0486" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0487" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0488" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0489" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0490" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0491" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0492" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0493" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0494" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0495" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0496" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0497" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0498" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0499" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0500" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0501" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0502" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0503" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0504" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0505" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0506" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0507" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0508" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0509" role="lcghm">
            <property role="lacIc" value="	return u, err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0510" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0511" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0512" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0513" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0514" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0515" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0516" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0517" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0518" role="3cqZAp" />
        <!-- // List -->
        <node concept="lc7rE" id="#Ia_0519" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0520" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0521" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0522" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0523" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0524" role="lcghm">
            <property role="lacIc" value=") List() ([]" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0525" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0526" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0527" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0528" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0529" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0530" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0531" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0532" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0533" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0534" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0535" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0536" role="lcghm">
            <property role="lacIc" value="		`SELECT id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0537" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0538" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0539" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0540" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0541" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0542" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0543" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0544" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0545" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0546" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0547" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0548" role="lcghm">
            <property role="lacIc" value=" ORDER BY id`)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0549" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0550" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0551" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0552" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0553" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0554" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0555" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0556" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0557" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0558" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0559" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0560" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0561" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0562" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0563" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0564" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0565" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0566" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0567" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0568" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0569" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0570" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0571" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0572" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0573" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0574" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0575" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0576" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0577" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0578" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0579" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0580" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0581" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0582" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0583" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0584" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;u.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0585" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0586" role="lcghm">
            <property role="lacIc" value=", &amp;u." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0587" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0588" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0589" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0590" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0591" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0592" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0593" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0594" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0595" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0596" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0597" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0598" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0599" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0600" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0601" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0602" role="lcghm">
            <property role="lacIc" value="		items = append(items, u)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0603" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0604" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0605" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0606" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0607" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0608" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0609" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0610" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0611" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0612" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0613" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0614" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0615" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0616" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0617" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0618" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0619" role="3cqZAp" />
        <!-- // Update -->
        <node concept="lc7rE" id="#Ia_0620" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0621" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0622" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0623" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0624" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0625" role="lcghm">
            <property role="lacIc" value=") Update(u *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0626" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0627" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0628" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0629" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0630" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0631" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0632" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0633" role="lcghm">
            <property role="lacIc" value="	return r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0634" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0635" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0636" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0637" role="lcghm">
            <property role="lacIc" value="		`UPDATE " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0638" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0639" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0640" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0641" role="lcghm">
            <property role="lacIc" value=" SET " />
          </node>
        </node>
        <!-- VAR: int uidx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- ─── IF: if (f.dataType != DataType.timestamp) ─── -->
        <!-- IF: if (uidx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0642" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0643" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: uidx = uidx + 1; -->
        <node concept="lc7rE" id="#Ia_0644" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0645" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0646" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0647" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0648" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0649" role="lcghm">
            <property role="lacIc" value="▶uidx◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- VAR: int nextParam = uidx + 1; -->
        <node concept="lc7rE" id="#Ia_0650" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0651" role="lcghm">
            <property role="lacIc" value=", updated_at = NOW()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0652" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0653" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0654" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0655" role="lcghm">
            <property role="lacIc" value="		 WHERE id = $" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0656" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0657" role="lcghm">
            <property role="lacIc" value="▶nextParam◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0658" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0659" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0660" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0661" role="lcghm">
            <property role="lacIc" value="		 RETURNING updated_at`," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0662" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0663" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <!-- IF: if (f.dataType != DataType.timestamp) -->
        <node concept="lc7rE" id="#Ia_0664" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0665" role="lcghm">
            <property role="lacIc" value="		u." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0666" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0667" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0668" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0669" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0670" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0671" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0672" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0673" role="lcghm">
            <property role="lacIc" value="		u.ID," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0674" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0675" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0676" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0677" role="lcghm">
            <property role="lacIc" value="	).Scan(&amp;u.UpdatedAt)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0678" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0679" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0680" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0681" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0682" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0683" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0684" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0685" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0686" role="3cqZAp" />
        <!-- // Delete -->
        <node concept="lc7rE" id="#Ia_0687" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0688" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0689" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0690" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0691" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0692" role="lcghm">
            <property role="lacIc" value=") Delete(id int64) error {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0693" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0694" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0695" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0696" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0697" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0698" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0699" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0700" role="lcghm">
            <property role="lacIc" value=" WHERE id = $1`, id)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0701" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0702" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0703" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0704" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0705" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0706" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0707" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0708" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0709" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0710" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0711" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0712" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_0713" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Repositories — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="#Is_0714" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="3clFbH" id="#Is_0715" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_0716" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0717" role="lcghm">
            <property role="lacIc" value="type " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0718" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0719" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0720" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0721" role="lcghm">
            <property role="lacIc" value=" struct{ db *sql.DB }" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0722" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0723" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0724" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0725" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0726" role="3cqZAp" />
        <!-- // Assign -->
        <node concept="lc7rE" id="#Ia_0727" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0728" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0729" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0730" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0731" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0732" role="lcghm">
            <property role="lacIc" value=") Assign(" />
          </node>
        </node>
        <!-- VAR: int argIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (argIdx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0733" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0734" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="#Ia_0735" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0736" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0737" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0738" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <!-- ASSIGN: argIdx = argIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- // query -->
        <node concept="lc7rE" id="#Ia_0739" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0740" role="lcghm">
            <property role="lacIc" value=") (*" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0741" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0742" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0743" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0744" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0745" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0746" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0747" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0748" role="lcghm">
            <property role="lacIc" value="	ur := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0749" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0750" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0751" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0752" role="lcghm">
            <property role="lacIc" value="{}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0753" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0754" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0755" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0756" role="lcghm">
            <property role="lacIc" value="	err := r.db.QueryRow(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0757" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0758" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0759" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0760" role="lcghm">
            <property role="lacIc" value="		`INSERT INTO " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0761" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0762" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0763" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0764" role="lcghm">
            <property role="lacIc" value=" (" />
          </node>
        </node>
        <!-- VAR: int fkIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (fkIdx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0765" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0766" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="#Ia_0767" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0768" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- ASSIGN: fkIdx = fkIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0769" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0770" role="lcghm">
            <property role="lacIc" value=") VALUES (" />
          </node>
        </node>
        <!-- ═══ FOR: for (int i = 1; i &lt;= fkIdx; i++) ═══ -->
        <!-- IF: if (i &gt; 1) -->
        <node concept="lc7rE" id="#Ia_0771" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0772" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="#Ia_0773" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0774" role="lcghm">
            <property role="lacIc" value="$" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0775" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0776" role="lcghm">
            <property role="lacIc" value="▶i◀" />
          </node>
        </node>
        <!-- END FOR -->
        <node concept="lc7rE" id="#Ia_0777" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0778" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0779" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0780" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0781" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0782" role="lcghm">
            <property role="lacIc" value="		 ON CONFLICT (" />
          </node>
        </node>
        <!-- VAR: int ckIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (ckIdx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0783" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0784" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="#Ia_0785" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0786" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- ASSIGN: ckIdx = ckIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0787" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0788" role="lcghm">
            <property role="lacIc" value=") DO NOTHING" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0789" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0790" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0791" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0792" role="lcghm">
            <property role="lacIc" value="		 RETURNING " />
          </node>
        </node>
        <!-- VAR: int retIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- IF: if (retIdx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0793" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0794" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="#Ia_0795" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0796" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0797" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0798" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: retIdx = retIdx + 1; -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0799" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0800" role="lcghm">
            <property role="lacIc" value="`," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0801" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0802" role="lcghm" />
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="#Ia_0803" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0804" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0805" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0806" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0807" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0808" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0809" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0810" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0811" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0812" role="lcghm">
            <property role="lacIc" value="	).Scan(" />
          </node>
        </node>
        <!-- VAR: int scanIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- IF: if (scanIdx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0813" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0814" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="#Ia_0815" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0816" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0817" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0818" role="lcghm">
            <property role="lacIc" value="▶fr.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ─── IF: if (field.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) field; -->
        <node concept="lc7rE" id="#Ia_0819" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0820" role="lcghm">
            <property role="lacIc" value="&amp;ur." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0821" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0822" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: scanIdx = scanIdx + 1; -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0823" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0824" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0825" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0826" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0827" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0828" role="lcghm">
            <property role="lacIc" value="	return ur, err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0829" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0830" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0831" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0832" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0833" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0834" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0835" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0836" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0837" role="3cqZAp" />
        <!-- // Remove -->
        <node concept="lc7rE" id="#Ia_0838" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0839" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0840" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0841" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0842" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0843" role="lcghm">
            <property role="lacIc" value=") Remove(" />
          </node>
        </node>
        <!-- VAR: int rmIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (rmIdx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0844" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0845" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <!-- END IF -->
        <node concept="lc7rE" id="#Ia_0846" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0847" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0848" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0849" role="lcghm">
            <property role="lacIc" value=" int64" />
          </node>
        </node>
        <!-- ASSIGN: rmIdx = rmIdx + 1; -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0850" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0851" role="lcghm">
            <property role="lacIc" value=") error {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0852" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0853" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0854" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0855" role="lcghm">
            <property role="lacIc" value="	_, err := r.db.Exec(`DELETE FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0856" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0857" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0858" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0859" role="lcghm">
            <property role="lacIc" value=" WHERE " />
          </node>
        </node>
        <!-- VAR: int wIdx = 0; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (wIdx &gt; 0) -->
        <node concept="lc7rE" id="#Ia_0860" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0861" role="lcghm">
            <property role="lacIc" value=" AND " />
          </node>
        </node>
        <!-- END IF -->
        <!-- ASSIGN: wIdx = wIdx + 1; -->
        <node concept="lc7rE" id="#Ia_0862" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0863" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0864" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0865" role="lcghm">
            <property role="lacIc" value=" = $" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0866" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0867" role="lcghm">
            <property role="lacIc" value="▶wIdx◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0868" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0869" role="lcghm">
            <property role="lacIc" value="`" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <node concept="lc7rE" id="#Ia_0870" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0871" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0872" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0873" role="lcghm">
            <property role="lacIc" value="▶fr.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0874" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0875" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0876" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0877" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0878" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0879" role="lcghm">
            <property role="lacIc" value="	return err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0880" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0881" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0882" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0883" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0884" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0885" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0886" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0887" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_0888" role="3cqZAp" />
        <!-- // Cross-queries -->
        <!-- ═══ FOREACH: foreach refA in schema.Fields ═══ -->
        <!-- ─── IF: if (refA.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA; -->
        <!-- ═══ FOREACH: foreach refB in schema.Fields ═══ -->
        <!-- ─── IF: if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB; -->
        <!-- VAR: string ts = frB.target_schema.structName(); -->
        <!-- VAR: string tt = frB.target_schema.name; -->
        <node concept="3clFbH" id="#Is_0889" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_0890" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0891" role="lcghm">
            <property role="lacIc" value="func (r *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0892" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0893" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0894" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0895" role="lcghm">
            <property role="lacIc" value=") Get" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0896" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0897" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0898" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0899" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0900" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0901" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0902" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0903" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0904" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0905" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0906" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0907" role="lcghm">
            <property role="lacIc" value=" int64) ([]" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0908" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0909" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0910" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0911" role="lcghm">
            <property role="lacIc" value=", error) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0912" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0913" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0914" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0915" role="lcghm">
            <property role="lacIc" value="	rows, err := r.db.Query(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0916" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0917" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0918" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0919" role="lcghm">
            <property role="lacIc" value="		`SELECT t.id" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach tf in frB.target_schema.Fields ═══ -->
        <!-- ─── IF: if (tf.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) tf; -->
        <node concept="lc7rE" id="#Ia_0920" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0921" role="lcghm">
            <property role="lacIc" value=", t." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0922" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0923" role="lcghm">
            <property role="lacIc" value="▶f.name◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_0924" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0925" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0926" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0927" role="lcghm">
            <property role="lacIc" value="		 FROM " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0928" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0929" role="lcghm">
            <property role="lacIc" value="▶tt◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0930" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0931" role="lcghm">
            <property role="lacIc" value=" t" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0932" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0933" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0934" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0935" role="lcghm">
            <property role="lacIc" value="		 INNER JOIN " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0936" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0937" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0938" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0939" role="lcghm">
            <property role="lacIc" value=" j ON j." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0940" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0941" role="lcghm">
            <property role="lacIc" value="▶frB.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0942" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0943" role="lcghm">
            <property role="lacIc" value=" = t.id" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0944" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0945" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0946" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0947" role="lcghm">
            <property role="lacIc" value="		 WHERE j." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0948" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0949" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0950" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0951" role="lcghm">
            <property role="lacIc" value=" = $1" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0952" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0953" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0954" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0955" role="lcghm">
            <property role="lacIc" value="		 ORDER BY t.id`, " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0956" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0957" role="lcghm">
            <property role="lacIc" value="▶frA.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0958" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0959" role="lcghm">
            <property role="lacIc" value="," />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0960" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0961" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0962" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0963" role="lcghm">
            <property role="lacIc" value="	)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0964" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0965" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0966" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0967" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0968" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0969" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0970" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0971" role="lcghm">
            <property role="lacIc" value="		return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0972" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0973" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0974" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0975" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0976" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0977" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0978" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0979" role="lcghm">
            <property role="lacIc" value="	defer rows.Close()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0980" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0981" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0982" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0983" role="lcghm">
            <property role="lacIc" value="	var items []" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0984" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0985" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0986" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0987" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0988" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0989" role="lcghm">
            <property role="lacIc" value="	for rows.Next() {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0990" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0991" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0992" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0993" role="lcghm">
            <property role="lacIc" value="		var item " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0994" role="3cqZAp">
          <node concept="la8eA" id="#Ix_0995" role="lcghm">
            <property role="lacIc" value="▶ts◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_0996" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_0997" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_0998" role="3cqZAp">
          <node concept="la8eA" id="#Ic_0999" role="lcghm">
            <property role="lacIc" value="		if err := rows.Scan(&amp;item.ID" />
          </node>
        </node>
        <!-- ═══ FOREACH: foreach tf in frB.target_schema.Fields ═══ -->
        <!-- ─── IF: if (tf.isInstanceOf(Field)) ─── -->
        <!-- VAR: node&lt;Field&gt; f = (Field) tf; -->
        <node concept="lc7rE" id="#Ia_1000" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1001" role="lcghm">
            <property role="lacIc" value=", &amp;item." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1002" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1003" role="lcghm">
            <property role="lacIc" value="▶f.pascalName()◀" />
          </node>
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="lc7rE" id="#Ia_1004" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1005" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1006" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1007" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1008" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1009" role="lcghm">
            <property role="lacIc" value="			return nil, err" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1010" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1011" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1012" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1013" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1014" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1015" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1016" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1017" role="lcghm">
            <property role="lacIc" value="		items = append(items, item)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1018" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1019" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1020" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1021" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1022" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1023" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1024" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1025" role="lcghm">
            <property role="lacIc" value="	return items, rows.Err()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1026" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1027" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1028" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1029" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1030" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1031" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1032" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1033" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_1034" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — regular schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="#Is_1035" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <node concept="3clFbH" id="#Is_1036" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_1037" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1038" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1039" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1040" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1041" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1042" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1043" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1044" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1045" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1046" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1047" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1048" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1049" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1050" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1051" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1052" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1053" role="3cqZAp" />
        <!-- // Create handler -->
        <node concept="lc7rE" id="#Ia_1054" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1055" role="lcghm">
            <property role="lacIc" value="func handleCreate" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1056" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1057" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1058" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1059" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1060" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1061" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1062" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1063" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1064" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1065" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1066" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1067" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1068" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1069" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1070" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1071" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1072" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1073" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1074" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1075" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1076" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1077" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1078" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1079" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1080" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1081" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1082" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1083" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1084" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1085" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1086" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1087" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1088" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1089" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1090" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1091" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1092" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1093" role="lcghm">
            <property role="lacIc" value="		if err := repo.Create(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1094" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1095" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1096" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1097" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1098" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1099" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1100" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1101" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1102" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1103" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1104" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1105" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1106" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1107" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1108" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1109" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1110" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1111" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1112" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1113" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1114" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1115" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1116" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1117" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1118" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1119" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1120" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1121" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1122" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1123" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1124" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1125" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1126" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1127" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1128" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1129" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1130" role="3cqZAp" />
        <!-- // Get handler -->
        <node concept="lc7rE" id="#Ia_1131" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1132" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1133" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1134" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1135" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1136" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1137" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1138" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1139" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1140" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1141" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1142" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1143" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1144" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1145" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1146" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1147" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1148" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1149" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1150" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1151" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1152" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1153" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1154" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1155" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1156" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1157" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1158" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1159" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1160" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1161" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1162" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1163" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1164" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1165" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1166" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1167" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1168" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1169" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1170" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1171" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1172" role="lcghm">
            <property role="lacIc" value="		u, err := repo.GetByID(id)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1173" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1174" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1175" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1176" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1177" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1178" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1179" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1180" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusNotFound)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1181" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1182" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1183" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1184" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1185" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1186" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1187" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1188" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1189" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1190" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1191" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1192" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1193" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1194" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1195" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1196" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1197" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1198" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1199" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1200" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1201" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1202" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1203" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1204" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1205" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1206" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1207" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1208" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1209" role="3cqZAp" />
        <!-- // List handler -->
        <node concept="lc7rE" id="#Ia_1210" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1211" role="lcghm">
            <property role="lacIc" value="func handleList" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1212" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1213" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1214" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1215" role="lcghm">
            <property role="lacIc" value="s(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1216" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1217" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1218" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1219" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1220" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1221" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1222" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1223" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1224" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1225" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1226" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1227" role="lcghm">
            <property role="lacIc" value="		items, err := repo.List()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1228" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1229" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1230" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1231" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1232" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1233" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1234" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1235" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1236" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1237" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1238" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1239" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1240" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1241" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1242" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1243" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1244" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1245" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1246" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1247" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1248" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1249" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1250" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1251" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1252" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1253" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1254" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1255" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1256" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1257" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1258" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1259" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1260" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1261" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1262" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1263" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1264" role="3cqZAp" />
        <!-- // Update handler -->
        <node concept="lc7rE" id="#Ia_1265" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1266" role="lcghm">
            <property role="lacIc" value="func handleUpdate" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1267" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1268" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1269" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1270" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1271" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1272" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1273" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1274" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1275" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1276" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1277" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1278" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1279" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1280" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1281" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1282" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1283" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1284" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1285" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1286" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1287" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1288" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1289" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1290" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1291" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1292" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1293" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1294" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1295" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1296" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1297" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1298" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1299" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1300" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1301" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1302" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1303" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1304" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1305" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1306" role="lcghm">
            <property role="lacIc" value="		var u " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1307" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1308" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1309" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1310" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1311" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1312" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1313" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1314" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1315" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1316" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1317" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1318" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1319" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1320" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1321" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1322" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1323" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1324" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1325" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1326" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1327" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1328" role="lcghm">
            <property role="lacIc" value="		u.ID = id" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1329" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1330" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1331" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1332" role="lcghm">
            <property role="lacIc" value="		if err := repo.Update(&amp;u); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1333" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1334" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1335" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1336" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1337" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1338" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1339" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1340" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1341" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1342" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1343" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1344" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1345" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1346" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1347" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1348" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1349" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1350" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1351" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1352" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(u)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1353" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1354" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1355" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1356" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1357" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1358" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1359" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1360" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1361" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1362" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1363" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1364" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1365" role="3cqZAp" />
        <!-- // Delete handler -->
        <node concept="lc7rE" id="#Ia_1366" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1367" role="lcghm">
            <property role="lacIc" value="func handleDelete" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1368" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1369" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1370" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1371" role="lcghm">
            <property role="lacIc" value="(repo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1372" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1373" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1374" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1375" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1376" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1377" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1378" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1379" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1380" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1381" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1382" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1383" role="lcghm">
            <property role="lacIc" value="		idStr := r.PathValue(&quot;id&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1384" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1385" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1386" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1387" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(idStr, 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1388" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1389" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1390" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1391" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1392" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1393" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1394" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1395" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1396" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1397" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1398" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1399" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1400" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1401" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1402" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1403" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1404" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1405" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1406" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1407" role="lcghm">
            <property role="lacIc" value="		if err := repo.Delete(id); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1408" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1409" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1410" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1411" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1412" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1413" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1414" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1415" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1416" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1417" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1418" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1419" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1420" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1421" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1422" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1423" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1424" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1425" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1426" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1427" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1428" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1429" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1430" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1431" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1432" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1433" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1434" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1435" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_1436" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // HTTP Handlers — join schemas -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="#Is_1437" role="3cqZAp" />
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string rn = schema.repoName(); -->
        <node concept="3clFbH" id="#Is_1438" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_1439" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1440" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1441" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1442" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1443" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1444" role="lcghm">
            <property role="lacIc" value="// HTTP Handlers — " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1445" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1446" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1447" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1448" role="lcghm">
            <property role="lacIc" value=" (assignments)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1449" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1450" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1451" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1452" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1453" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1454" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1455" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1456" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1457" role="3cqZAp" />
        <!-- VAR: node&lt;FieldRefrence&gt; firstRef = null; -->
        <!-- VAR: node&lt;FieldRefrence&gt; secondRef = null; -->
        <!-- ═══ FOREACH: foreach field in schema.Fields ═══ -->
        <!-- ─── IF: if (field.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; fr = (FieldRefrence) field; -->
        <!-- IF: if (firstRef == null) { firstRef = fr; } -->
        <!-- IF: else if (secondRef == null) { secondRef = fr; } -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_1458" role="3cqZAp" />
        <!-- // Assign handler -->
        <node concept="lc7rE" id="#Ia_1459" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1460" role="lcghm">
            <property role="lacIc" value="func handleAssign" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1461" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1462" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1463" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1464" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1465" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1466" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1467" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1468" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1469" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1470" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1471" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1472" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1473" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1474" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1475" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1476" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1477" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1478" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1479" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1480" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1481" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1482" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1483" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1484" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1485" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1486" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1487" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1488" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1489" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1490" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1491" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1492" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1493" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1494" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1495" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1496" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1497" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1498" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1499" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1500" role="lcghm">
            <property role="lacIc" value="		var body Assign" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1501" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1502" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1503" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1504" role="lcghm">
            <property role="lacIc" value="Body" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1505" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1506" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1507" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1508" role="lcghm">
            <property role="lacIc" value="		if err := json.NewDecoder(r.Body).Decode(&amp;body); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1509" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1510" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1511" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1512" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1513" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1514" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1515" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1516" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1517" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1518" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1519" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1520" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1521" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1522" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1523" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1524" role="lcghm">
            <property role="lacIc" value="		ur, err := urRepo.Assign(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1525" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1526" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1527" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1528" role="lcghm">
            <property role="lacIc" value=", body." />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1529" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1530" role="lcghm">
            <property role="lacIc" value="▶secondRef.pascalName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1531" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1532" role="lcghm">
            <property role="lacIc" value=")" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1533" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1534" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1535" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1536" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1537" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1538" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1539" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1540" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1541" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1542" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1543" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1544" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1545" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1546" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1547" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1548" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1549" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1550" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1551" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1552" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1553" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1554" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1555" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1556" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusCreated)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1557" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1558" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1559" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1560" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(ur)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1561" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1562" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1563" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1564" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1565" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1566" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1567" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1568" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1569" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1570" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1571" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1572" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1573" role="3cqZAp" />
        <!-- // Remove handler -->
        <node concept="lc7rE" id="#Ia_1574" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1575" role="lcghm">
            <property role="lacIc" value="func handleRemove" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1576" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1577" role="lcghm">
            <property role="lacIc" value="▶secondRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1578" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1579" role="lcghm">
            <property role="lacIc" value="(urRepo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1580" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1581" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1582" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1583" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1584" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1585" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1586" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1587" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1588" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1589" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1590" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1591" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1592" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1593" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1594" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1595" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1596" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1597" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1598" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1599" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1600" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1601" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1602" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1603" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1604" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1605" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1606" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1607" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1608" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1609" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1610" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1611" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1612" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1613" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1614" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1615" role="lcghm">
            <property role="lacIc" value="		" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1616" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1617" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1618" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1619" role="lcghm">
            <property role="lacIc" value=", err := strconv.ParseInt(r.PathValue(&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1620" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1621" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1622" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1623" role="lcghm">
            <property role="lacIc" value="&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1624" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1625" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1626" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1627" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1628" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1629" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1630" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1631" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1632" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1633" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1634" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1635" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1636" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1637" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1638" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1639" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1640" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1641" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1642" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1643" role="lcghm">
            <property role="lacIc" value="		if err := urRepo.Remove(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1644" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1645" role="lcghm">
            <property role="lacIc" value="▶firstRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1646" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1647" role="lcghm">
            <property role="lacIc" value=", " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1648" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1649" role="lcghm">
            <property role="lacIc" value="▶secondRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1650" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1651" role="lcghm">
            <property role="lacIc" value="); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1652" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1653" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1654" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1655" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1656" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1657" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1658" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1659" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1660" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1661" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1662" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1663" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1664" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1665" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1666" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1667" role="lcghm">
            <property role="lacIc" value="		w.WriteHeader(http.StatusNoContent)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1668" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1669" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1670" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1671" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1672" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1673" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1674" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1675" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1676" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1677" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1678" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1679" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1680" role="3cqZAp" />
        <!-- // Cross-query handlers -->
        <!-- ═══ FOREACH: foreach refA in schema.Fields ═══ -->
        <!-- ─── IF: if (refA.isInstanceOf(FieldRefrence)) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frA = (FieldRefrence) refA; -->
        <!-- ═══ FOREACH: foreach refB in schema.Fields ═══ -->
        <!-- ─── IF: if (refB.isInstanceOf(FieldRefrence) &amp;&amp; refB != refA) ─── -->
        <!-- VAR: node&lt;FieldRefrence&gt; frB = (FieldRefrence) refB; -->
        <node concept="lc7rE" id="#Ia_1681" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1682" role="lcghm">
            <property role="lacIc" value="func handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1683" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1684" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1685" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1686" role="lcghm">
            <property role="lacIc" value="▶frB.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1687" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1688" role="lcghm">
            <property role="lacIc" value="s(urRepo *" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1689" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1690" role="lcghm">
            <property role="lacIc" value="▶rn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1691" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1692" role="lcghm">
            <property role="lacIc" value=") http.HandlerFunc {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1693" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1694" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1695" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1696" role="lcghm">
            <property role="lacIc" value="	return func(w http.ResponseWriter, r *http.Request) {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1697" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1698" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1699" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1700" role="lcghm">
            <property role="lacIc" value="		id, err := strconv.ParseInt(r.PathValue(&quot;id&quot;), 10, 64)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1701" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1702" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1703" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1704" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1705" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1706" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1707" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1708" role="lcghm">
            <property role="lacIc" value="			http.Error(w, &quot;invalid id&quot;, http.StatusBadRequest)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1709" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1710" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1711" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1712" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1713" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1714" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1715" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1716" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1717" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1718" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1719" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1720" role="lcghm">
            <property role="lacIc" value="		items, err := urRepo.Get" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1721" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1722" role="lcghm">
            <property role="lacIc" value="▶frB.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1723" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1724" role="lcghm">
            <property role="lacIc" value="sBy" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1725" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1726" role="lcghm">
            <property role="lacIc" value="▶frA.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1727" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1728" role="lcghm">
            <property role="lacIc" value="(id)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1729" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1730" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1731" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1732" role="lcghm">
            <property role="lacIc" value="		if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1733" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1734" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1735" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1736" role="lcghm">
            <property role="lacIc" value="			http.Error(w, err.Error(), http.StatusInternalServerError)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1737" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1738" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1739" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1740" role="lcghm">
            <property role="lacIc" value="			return" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1741" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1742" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1743" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1744" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1745" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1746" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1747" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1748" role="lcghm">
            <property role="lacIc" value="		w.Header().Set(&quot;Content-Type&quot;, &quot;application/json&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1749" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1750" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1751" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1752" role="lcghm">
            <property role="lacIc" value="		json.NewEncoder(w).Encode(items)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1753" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1754" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1755" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1756" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1757" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1758" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1759" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1760" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1761" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1762" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1763" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1764" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_1765" role="3cqZAp" />
        <!-- // ============================================================ -->
        <!-- // Main -->
        <!-- // ============================================================ -->
        <node concept="3clFbH" id="#Is_1766" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_1767" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1768" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1769" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1770" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1771" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1772" role="lcghm">
            <property role="lacIc" value="// Main" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1773" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1774" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1775" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1776" role="lcghm">
            <property role="lacIc" value="// ============================================================" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1777" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1778" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1779" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1780" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1781" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1782" role="lcghm">
            <property role="lacIc" value="func main() {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1783" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1784" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1785" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1786" role="lcghm">
            <property role="lacIc" value="	dbURL := os.Getenv(&quot;DATABASE_URL&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1787" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1788" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1789" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1790" role="lcghm">
            <property role="lacIc" value="	if dbURL == &quot;&quot; {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1791" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1792" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1793" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1794" role="lcghm">
            <property role="lacIc" value="		dbURL = &quot;postgres://" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1795" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1796" role="lcghm">
            <property role="lacIc" value="▶infra.dbUser◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1797" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1798" role="lcghm">
            <property role="lacIc" value=":" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1799" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1800" role="lcghm">
            <property role="lacIc" value="▶infra.dbPass◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1801" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1802" role="lcghm">
            <property role="lacIc" value="@localhost:5432/" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1803" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1804" role="lcghm">
            <property role="lacIc" value="▶infra.dbName◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1805" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1806" role="lcghm">
            <property role="lacIc" value="?sslmode=disable&quot;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1807" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1808" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1809" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1810" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1811" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1812" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1813" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1814" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1815" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1816" role="lcghm">
            <property role="lacIc" value="	db, err := sql.Open(&quot;postgres&quot;, dbURL)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1817" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1818" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1819" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1820" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1821" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1822" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1823" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1824" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1825" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1826" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1827" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1828" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1829" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1830" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1831" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1832" role="lcghm">
            <property role="lacIc" value="	defer db.Close()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1833" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1834" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1835" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1836" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1837" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1838" role="lcghm">
            <property role="lacIc" value="	for i := 0; i &lt; 5; i++ {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1839" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1840" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1841" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1842" role="lcghm">
            <property role="lacIc" value="		if err = db.Ping(); err == nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1843" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1844" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1845" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1846" role="lcghm">
            <property role="lacIc" value="			break" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1847" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1848" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1849" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1850" role="lcghm">
            <property role="lacIc" value="		}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1851" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1852" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1853" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1854" role="lcghm">
            <property role="lacIc" value="		log.Printf(&quot;DB not ready, retrying... (%d/5)&quot;, i+1)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1855" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1856" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1857" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1858" role="lcghm">
            <property role="lacIc" value="		time.Sleep(2 * time.Second)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1859" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1860" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1861" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1862" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1863" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1864" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1865" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1866" role="lcghm">
            <property role="lacIc" value="	if err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1867" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1868" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1869" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1870" role="lcghm">
            <property role="lacIc" value="		log.Fatal(&quot;DB connection failed:&quot;, err)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1871" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1872" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1873" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1874" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1875" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1876" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1877" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1878" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1879" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1880" role="lcghm">
            <property role="lacIc" value="	if _, err := db.Exec(migrationSQL); err != nil {" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1881" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1882" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1883" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1884" role="lcghm">
            <property role="lacIc" value="		log.Fatal(err)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1885" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1886" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1887" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1888" role="lcghm">
            <property role="lacIc" value="	}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1889" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1890" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1891" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1892" role="lcghm">
            <property role="lacIc" value="	log.Println(&quot;Migration complete&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1893" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1894" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1895" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1896" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1897" role="3cqZAp" />
        <!-- // Repo instantiation -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <node concept="lc7rE" id="#Ia_1898" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1899" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1900" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1901" role="lcghm">
            <property role="lacIc" value="▶schema.singularName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1902" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1903" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1904" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1905" role="lcghm">
            <property role="lacIc" value="▶schema.repoName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1906" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1907" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1908" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1909" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (schema.hasReferences()) ─── -->
        <node concept="lc7rE" id="#Ia_1910" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1911" role="lcghm">
            <property role="lacIc" value="	" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1912" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1913" role="lcghm">
            <property role="lacIc" value="▶schema.singularName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1914" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1915" role="lcghm">
            <property role="lacIc" value="Repo := &amp;" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1916" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1917" role="lcghm">
            <property role="lacIc" value="▶schema.repoName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1918" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1919" role="lcghm">
            <property role="lacIc" value="{db: db}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1920" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1921" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_1922" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_1923" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1924" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1925" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1926" role="lcghm">
            <property role="lacIc" value="	mux := http.NewServeMux()" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1927" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1928" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1929" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1930" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1931" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1932" role="lcghm">
            <property role="lacIc" value="	// Swagger UI" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1933" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1934" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1935" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1936" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /swagger/*&quot;, httpSwagger.WrapHandler)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1937" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1938" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1939" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1940" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_1941" role="3cqZAp" />
        <!-- // Regular routes -->
        <!-- ═══ FOREACH: foreach schema in model.schemas ═══ -->
        <!-- ─── IF: if (!(schema.hasReferences())) ─── -->
        <!-- VAR: string vr = schema.singularName() + "Repo"; -->
        <!-- VAR: string sn = schema.structName(); -->
        <!-- VAR: string tn = schema.name; -->
        <node concept="lc7rE" id="#Ia_1942" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1943" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1944" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1945" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1946" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1947" role="lcghm">
            <property role="lacIc" value="s" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1948" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1949" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1950" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1951" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1952" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1953" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1954" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1955" role="lcghm">
            <property role="lacIc" value="&quot;, handleCreate" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1956" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1957" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1958" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1959" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1960" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1961" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1962" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1963" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1964" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1965" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1966" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1967" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1968" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1969" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1970" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1971" role="lcghm">
            <property role="lacIc" value="&quot;, handleList" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1972" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1973" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1974" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1975" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1976" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1977" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1978" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1979" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1980" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1981" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1982" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1983" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1984" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1985" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1986" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1987" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1988" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1989" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1990" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1991" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1992" role="3cqZAp">
          <node concept="la8eA" id="#Ix_1993" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1994" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1995" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_1996" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_1997" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_1998" role="3cqZAp">
          <node concept="la8eA" id="#Ic_1999" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;PUT /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2000" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2001" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2002" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2003" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleUpdate" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2004" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2005" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2006" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2007" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2008" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2009" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2010" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2011" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2012" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2013" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2014" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2015" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2016" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2017" role="lcghm">
            <property role="lacIc" value="▶tn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2018" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2019" role="lcghm">
            <property role="lacIc" value="/{id}&quot;, handleDelete" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2020" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2021" role="lcghm">
            <property role="lacIc" value="▶sn◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2022" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2023" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2024" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2025" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2026" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2027" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2028" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2029" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2030" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2031" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_2032" role="3cqZAp" />
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
        <node concept="lc7rE" id="#Ia_2033" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2034" role="lcghm">
            <property role="lacIc" value="	// " />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2035" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2036" role="lcghm">
            <property role="lacIc" value="▶schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2037" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2038" role="lcghm">
            <property role="lacIc" value=" assignments" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2039" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2040" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2041" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2042" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;POST /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2043" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2044" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2045" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2046" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2047" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2048" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2049" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2050" role="lcghm">
            <property role="lacIc" value="&quot;, handleAssign" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2051" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2052" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2053" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2054" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2055" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2056" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2057" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2058" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2059" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2060" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2061" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2062" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;DELETE /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2063" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2064" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2065" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2066" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2067" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2068" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2069" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2070" role="lcghm">
            <property role="lacIc" value="/{" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2071" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2072" role="lcghm">
            <property role="lacIc" value="▶sRef.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2073" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2074" role="lcghm">
            <property role="lacIc" value="}&quot;, handleRemove" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2075" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2076" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2077" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2078" role="lcghm">
            <property role="lacIc" value="(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2079" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2080" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2081" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2082" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2083" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2084" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2085" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2086" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2087" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2088" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2089" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2090" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2091" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2092" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2093" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2094" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2095" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2096" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2097" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2098" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2099" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2100" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2101" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2102" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2103" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2104" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2105" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2106" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2107" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2108" role="lcghm">
            <property role="lacIc" value="	mux.HandleFunc(&quot;GET /" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2109" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2110" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2111" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2112" role="lcghm">
            <property role="lacIc" value="/{id}/" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2113" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2114" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.name◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2115" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2116" role="lcghm">
            <property role="lacIc" value="&quot;, handleGet" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2117" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2118" role="lcghm">
            <property role="lacIc" value="▶sRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2119" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2120" role="lcghm">
            <property role="lacIc" value="▶fRef.target_schema.structName()◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2121" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2122" role="lcghm">
            <property role="lacIc" value="s(" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2123" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2124" role="lcghm">
            <property role="lacIc" value="▶vr◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2125" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2126" role="lcghm">
            <property role="lacIc" value="))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2127" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2128" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2129" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2130" role="lcghm" />
        </node>
        <!-- END IF -->
        <!-- END FOREACH -->
        <node concept="3clFbH" id="#Is_2131" role="3cqZAp" />
        <node concept="lc7rE" id="#Ia_2132" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2133" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Serving on :" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2134" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2135" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2136" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2137" role="lcghm">
            <property role="lacIc" value="&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2138" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2139" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2140" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2141" role="lcghm">
            <property role="lacIc" value="	fmt.Println(&quot;Swagger UI: http://localhost:" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2142" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2143" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2144" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2145" role="lcghm">
            <property role="lacIc" value="/swagger/index.html&quot;)" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2146" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2147" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2148" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2149" role="lcghm">
            <property role="lacIc" value="	log.Fatal(http.ListenAndServe(&quot;:" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2150" role="3cqZAp">
          <node concept="la8eA" id="#Ix_2151" role="lcghm">
            <property role="lacIc" value="▶infra.port◀" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2152" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2153" role="lcghm">
            <property role="lacIc" value="&quot;, mux))" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2154" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2155" role="lcghm" />
        </node>
        <node concept="lc7rE" id="#Ia_2156" role="3cqZAp">
          <node concept="la8eA" id="#Ic_2157" role="lcghm">
            <property role="lacIc" value="}" />
          </node>
        </node>
        <node concept="lc7rE" id="#Ia_2158" role="3cqZAp">
          <node concept="l8MVK" id="#Inl_2159" role="lcghm" />
        </node>
        <node concept="3clFbH" id="#Is_2160" role="3cqZAp" />
        <!-- END ? -->
        <!-- END ? -->
      </node>
    </node>
  </node>
</model>
