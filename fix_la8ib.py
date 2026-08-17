import re

with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'r') as f:
    text = f.read()

# Fix lc7rE (AppendOperation) to include the parts child (la8Ib)
app_match = re.search(r'(<concept[^>]+index="lc7rE">)', text)
if app_match:
    pos = app_match.end()
    # 1237306079178 is the ID for the 'part' child role
    addition = '\n        <child id="1237306115446" name="part" index="la8Ib" />'
    if 'index="la8Ib"' not in text[pos:pos+200]:
        text = text[:pos] + addition + text[pos:]
        print('✅ Added la8Ib to lc7rE')

with open('/home/prem-modha/MPSProjects/UserManagement/languages/UserManagement/models/UserManagement.textGen.mps', 'w') as f:
    f.write(text)

# Check one final time
all_used = set(re.findall(r'role="([^"]+)"', text))
all_registered = set(re.findall(r'<child[^>]+index="([^"]+)"', text)) | \
                 set(re.findall(r'<reference[^>]+index="([^"]+)"', text)) | \
                 set(re.findall(r'<property[^>]+index="([^"]+)"', text))
print(f'Missing roles: {all_used - all_registered}')
