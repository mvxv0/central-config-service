# Central Config Service - KT1

Ovaj projekat predstavlja web servis za centralizovano upravljanje konfiguracijama.

## Tehnologije
- **Jezik:** Go (Golang) 1.26.2
- **Ruter:** Gorilla Mux
- **Skladištenje:** In-memory (mape podataka)
- **ID** uuid biblioteka

## Funkcionalnosti (KT1)
U okviru ove faze implementirane su sledeće funkcionalnosti:
- **Konfiguracije:** Kompletan CRUD (Kreiranje, čitanje, brisanje).
- **Grupe konfiguracija:** Kreiranje, čitanje i brisanje grupa.
- **Povezivanje:** Dodavanje postojećih konfiguracija u specifične grupe (DTO model).
- **Validacija:** Osnovna provera postojanja resursa pre izvršavanja operacija.
- **Verzionisanje** Dobavljanje, slanje i brisanje podataka preko imena i verzije

## Struktura Projekta
- `model/`: Definicije struktura podataka (`Config`, `ConfigGroup`, `ConfigDTO`) i interfejsi.
- `repositories/`: Implementacija in-memory skladišta.
- `services/`: Poslovna logika aplikacije.
- `handlers/`: HTTP handler-i za obradu zahteva.

## API Rute

### Konfiguracije
| `POST` | `/configs` | Kreiranje nove konfiguracije |
| `GET` | `/configs` | Dobavljanje svih konfiguracija |
| `DELETE` | `/configs/{name}/{version}` | Brisanje konfiguracije |
| `GET` | `/configs/{name}/{version}` | Dobavljanje odredjene konfiguracije |

### Grupe
| `POST` | `/groups` | Kreiranje nove grupe |
| `GET` | `/groups/{name}/{version}` | Dobavljanje specifične grupe |
| `GET` | `/groups` | Dobavljanje svih grupa |
| `POST` | `/groups/{name}/{version}/configs/link` | Povezivanje postojeće konfiguracije sa grupom |
| `POST` | `/groups/{name}/{version}/configs` | Povezivanje nove konfiguracije sa grupom | 
| `DELETE` | `/groups/{name}/{version}` | Brisanje grupe |
| `DELETE` | `/groups/{name}/{version}/configs/{configName}` | Brisanje konfiguracije iz grupe |


## Pokretanje aplikacije

### Lokalno (Go)
Potrebno je imati instaliran Go. Pokrenite sledeću komandu u terminalu:
```bash
go run main.go

go mod tidy ( automatski dodavanje biblioteka u go.sum)