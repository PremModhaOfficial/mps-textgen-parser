import re
with open("parse_textgen.py", "r") as f:
    code = f.read()

old_logic = """                # Look for a ConceptDeclaration with the given name
                m = re.search(r'<node concept="1TIwiD" id="([^"]+)".*?name="{}"'.format(re.escape(concept_name)), content)
                if not m: # fallback if attributes are ordered differently
                    m = re.search(r'<node.*?id="([^"]+)".*?name="{}"'.format(re.escape(concept_name)), content)
                if m:
                    concept_id = m.group(1)"""

new_logic = """                # Look for a ConceptDeclaration (1TIwiD) with the given name (property TrG5h)
                nodes = re.findall(r'<node concept="1TIwiD" id="([^"]+)">(.*?)</node>', content, re.DOTALL)
                for nid, inner in nodes:
                    if re.search(rf'<property role="TrG5h" value="{re.escape(concept_name)}"', inner):
                        concept_id = nid
                        break
                
                # Fallback for older MPS versions where name is an attribute
                if not concept_id:
                    m = re.search(r'<node concept="1TIwiD" id="([^"]+)".*?name="{}"'.format(re.escape(concept_name)), content)
                    if m:
                        concept_id = m.group(1)"""

code = code.replace(old_logic, new_logic)

with open("parse_textgen.py", "w") as f:
    f.write(code)
