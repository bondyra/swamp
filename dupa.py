from bs4 import BeautifulSoup
import re

pesel_regex = re.compile(".*[0-9]{11}.*")
dupa = re.compile(".*Lp. (.+)")
lokal_regex = re.compile(".+PROJEKTOWYM NR ([^,]+),.*")

with open("pokazWydruk.htm", "r") as f:
    cuntent = f.read()

soup = BeautifulSoup(cuntent, "html.parser")

body = soup.html.body
children = list(body.children)

one = children[9]
two = children[11]

tab = list(one.children)[5].tbody  # <- tu sa roszczenia
wpisy = list(tab.children)

ludzie = {}
peoples = []
lokal = ""
for w in wpisy:
    if "Nr podstawy wpisu" in w.text:
        if lokal:
            ludzie[lokal] = peoples
        peoples = []
        lokal = ""
    if pesel_regex.match(w.text):
        m = dupa.match(w.text)
        if m:
            peoples.append(m.groups()[0])
    if "Treść wpisu" in w.text:
        m = lokal_regex.match(w.text.replace("\n",""))
        if m:
            lokal = m.groups()[0]


fullnames = [vv for k,v in ludzie.items() for vv in v]
pesels = [vv.split(",")[-1].strip() for k,v in ludzie.items() for vv in v]
pesel_to_no = {x: len([y for y in pesels if y==x]) for x in set(e for e in pesels)}

[k for k,v in ludzie.items() if any("85090909077" in vv for vv in v)]

#68040215716