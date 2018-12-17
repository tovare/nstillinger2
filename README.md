
Antall stillinger
=================

I desember 2018 fjernes det gamle stillingssøket slik at jeg ikke lenger kan
skrape siden for å få løpende tall slik som i nstillinger

Lager et endepunkt [](https://tovare.com/api/stillinger) som leverer 

    {"stillinger":"18357","annonser":"10413","nye":"0"}

Beregnet for å ligge bak en https proxy som NGNIX og init.d 

    Usage of nstillinger
    -p string
            Hvilket portnummer/adresse (default ":8085")
    -prefix string
            Hvilket adresse ligger løsningen på (default "/api")

Dette er mot et midlertidig API, vi får se hva som kommer. Foreløpige endepunkter hos NAV:

Endepunkt


URL [](https://stillingsok.nav.no/api/search")

Parametere

* q=<søketreng>
* sort=updated
* published=now-1d

Paginering:

* from=indeks fra 
* to=indeks til

Tov Are Jacobsen, 2018
