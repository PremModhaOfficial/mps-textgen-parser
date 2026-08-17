import re

with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'r') as f:
    text = f.read()

# baseLanguage (f3061a53-9226-4cc5-a443-f952ceaf5816) needs 3cpWs6, 17QB3L
base_lang_match = re.search(r'<language id="f3061a53-9226-4cc5-a443-f952ceaf5816" name="jetbrains.mps.baseLanguage">', text)
if base_lang_match:
    pos = base_lang_match.end()
    additions = ""
    if 'index="3cpWs6"' not in text:
        additions += '\n      <concept id="1068581242878" name="jetbrains.mps.baseLanguage.structure.ReturnStatement" flags="nn" index="3cpWs6">\n        <child id="1068581517676" name="expression" index="3cqZAk" />\n      </concept>'
    if 'index="17QB3L"' not in text:
        additions += '\n      <concept id="1225271177708" name="jetbrains.mps.baseLanguage.structure.StringType" flags="in" index="17QB3L" />'
    text = text[:pos] + additions + text[pos:]

# smodel (7866978e-a0f0-4cc7-81bc-4d213d9375e1) needs 3TrEf2
smodel_match = re.search(r'<language id="7866978e-a0f0-4cc7-81bc-4d213d9375e1" name="jetbrains.mps.lang.smodel">', text)
if smodel_match:
    pos = smodel_match.end()
    additions = ""
    if 'index="3TrEf2"' not in text:
        additions += '\n      <concept id="1138056143562" name="jetbrains.mps.lang.smodel.structure.SLinkAccess" flags="nn" index="3TrEf2">\n        <reference id="1138056516764" name="link" index="3Tt5mk" />\n      </concept>'
    text = text[:pos] + additions + text[pos:]

# Check textGen (b83431fe-5c8f-40bc-8a36-65e25f4dd253) for 29tfMY
textgen_match = re.search(r'<language id="b83431fe-5c8f-40bc-8a36-65e25f4dd253" name="jetbrains.mps.lang.textGen">', text)
if textgen_match:
    pos = textgen_match.end()
    additions = ""
    if 'index="29tfMY"' not in text:
        additions += '\n      <concept id="1233920501193" name="jetbrains.mps.lang.textGen.structure.FilenameFunction" flags="ig" index="29tfMY" />'
    text = text[:pos] + additions + text[pos:]

with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'w') as f:
    f.write(text)

used = set(re.findall(r'<node[^\>]+concept="([^"]+)"', text))
registered = set(re.findall(r'<concept[^\>]+index="([^"]+)"', text))
missing = used - registered
print(f'Missing concepts after fix: {missing}')
