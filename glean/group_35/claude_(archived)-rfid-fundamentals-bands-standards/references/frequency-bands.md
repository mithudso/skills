# RFID Frequency Bands & Coupling — Deep Reference

Companion to `../SKILL.md`. Physics of each band, the near-field/far-field boundary, and the
trade-offs that drive band selection. All figures are grounded in the sources listed in
SKILL.md; uncertain or approximate values are marked.

## The near-field / far-field boundary

The boundary between the near-field (reactive/inductive) and far-field (radiative) regions of
an antenna is approximately **λ / 2π**, where λ is the wavelength. This single number explains
why LF/HF behave fundamentally differently from UHF/microwave.

| Band | Freq | Wavelength λ | λ/2π (near-field edge) | Operating regime |
|---|---|---|---|---|
| LF | 125 kHz | ~2,400 m | ~382 m | Deep near-field; coupling distance ≪ λ |
| LF | 134.2 kHz | ~2,235 m | ~356 m | Deep near-field |
| HF | 13.56 MHz | ~22.1 m | ~3.5 m | Near-field (practical range ≪ λ/2π) |
| UHF | ~915 MHz | ~0.33 m | ~5.2 cm | Far-field at typical read distances |
| Microwave | 2.45 GHz | ~0.122 m | ~1.9 cm | Far-field |

Because the LF/HF wavelength is far larger than any practical reader antenna or read distance,
those systems live entirely in the near (magnetic) field and use **inductive coupling**. UHF
and microwave read distances exceed λ/2π, so they operate in the far (radiating) field and use
**backscatter**.

## Inductive coupling (LF & HF, near-field)

- **Energy transfer.** Reader and tag are two coils sharing a magnetic field (a loosely
  coupled air-core transformer). A changing current in the reader coil induces a current in
  the tag coil — *mutual inductance*. The tag's resonant circuit (coil + capacitor tuned to
  the carrier) maximizes the harvested voltage.
- **Field decay.** In the near-field the magnetic field strength H falls as **1/d³**. Because
  available power scales with the square of field strength, harvested **power falls as ~1/d⁶**.
  This steep roll-off is why inductive range is short and drops sharply with distance.
- **Tag → reader (load modulation).** The tag does not transmit its own carrier. It switches a
  load (resistor/capacitor) across its coil on and off in time with its data. This changes the
  reflected impedance seen by the reader and creates detectable **sidebands** around the
  carrier. The reader demodulates those sidebands.
- **Field strength requirements.** Proximity (ISO/IEC 14443) operates at **1.5–7.5 A/m**;
  vicinity (ISO/IEC 15693) operates at a weaker **0.15–5 A/m**, which is what lets vicinity
  cards work at greater distance with a lower-power field.
- **Material behaviour.** Magnetic near-fields pass through water, human tissue, wood, and
  most non-metals with little loss, so LF/HF read well around liquids and bodies. **Metal**
  near the coil detunes/absorbs the field (eddy currents), degrading reads. Orientation
  matters: the tag coil must cut enough flux, so coplanar/coaxial alignment reads best.

## Backscatter coupling (UHF & microwave, far-field)

- **Energy transfer.** The reader radiates an electromagnetic wave. The tag's antenna captures
  part of it; a rectifier (charge pump) converts RF to DC to power the chip. Modern passive UHF
  ICs are sensitive to roughly **−20 to −22 dBm or better** (approximate; vendor- and
  generation-dependent) — the lower the sensitivity figure, the longer the range.
- **Field decay.** Far-field power falls as **1/d²** (much gentler than near-field 1/d⁶),
  enabling multi-metre range. Note the *round-trip* penalty: the powering wave attenuates on
  the way out and the backscattered signal attenuates on the way back, so the return link is
  often the limiting factor.
- **Tag → reader (modulated backscatter).** The tag toggles a load across its antenna, which
  changes the antenna's reflection cross-section — making the tag alternately a good or poor
  reflector. The reader's directional coupler separates this weak reflected, modulated signal
  from its own outbound carrier. A tag antenna reflects strongly when its largest dimension is
  ≥ ~½ wavelength (≈16 cm at 915 MHz — hence typical UHF label dipole sizes).
- **Material behaviour.** UHF is **absorbed by water and tissue** (the 915 MHz band sits near
  water-absorption-relevant frequencies) and **detuned/reflected by metal**, so liquid-filled
  or metal items need specialty on-metal/liquid-tolerant tag designs and careful antenna
  placement. Microwave (2.45 GHz) is even more strongly absorbed by water.

## Band-by-band detail

### LF — 125 kHz / 134.2 kHz
- Coupling: inductive (near-field). Power class: typically passive.
- Range: a few cm up to ~10 cm. Data rate: lowest of all bands. Memory: small.
- Strengths: best penetration of liquid/tissue/non-metals; tolerant of orientation; works near
  the body. Weaknesses: very short range, slow, larger/coiled antennas, tags often cost more
  than HF.
- Uses: animal identification (ISO 11784/11785, ISO/IEC 14223), vehicle immobilizers, legacy
  proximity access cards (e.g. 125 kHz prox), some POS.

### HF — 13.56 MHz
- Coupling: inductive (near-field). Power class: passive. Globally available ISM band.
- Range: up to ~1 m for vicinity (15693); ~10 cm for proximity (14443) and NFC. Data rate:
  moderate — 106/212/424/848 kbit/s on ISO/IEC 14443; 15693 is slower (~26.48 kbit/s tag→reader
  high rate, down to ~1.65 kbit/s in 1-of-256 mode).
- Strengths: cheapest tags of the three bands; mature standards; good for short, controlled,
  tap-style reads; strong security ecosystem at the card-product layer. Weaknesses: short
  range; metal-sensitive.
- Uses: access/ID cards, contactless payment (EMVCo over 14443), transit, library media,
  e-passports, and **all NFC**.

### UHF — 860–960 MHz (passive); 433 MHz (active)
- Coupling: backscatter (far-field). Power class: passive, semi-passive (BAP), or active.
- Range: ~2–10 m passive (to ~15–20 ft with good antennas/tags); active 433 MHz reaches tens
  to hundreds of metres. Data rate: high — singulates **hundreds of tags per second**.
- Strengths: long range, fast bulk reading, low and falling tag cost, large supply-chain
  ecosystem. Weaknesses: absorbed by water/tissue, detuned by metal, **region-specific
  spectrum**, more multipath/nulls.
- Uses: supply chain, retail inventory, warehouse/logistics, apparel, returnable assets,
  toll/yard (often with active/BAP tags).

### Microwave — 2.45 GHz (also 5.8 GHz)
- Coupling: backscatter (far-field). Power class: typically active or semi-passive.
- Range: ~1–2 m for semi-passive; much further for active. Data rate: high.
- Strengths: small antennas, high data rate, fine spatial resolution. Weaknesses: strongest
  water absorption, costlier, most metal/multipath sensitivity.
- Uses: electronic toll collection (some systems), real-time location systems (RTLS), specific
  industrial/automotive applications. (RTLS application depth → asset-tracking sibling.)

## UHF regional spectrum (why "UHF" is not one band)

Unlike the global 13.56 MHz HF ISM band, UHF allocations differ by jurisdiction; a tag reads
across the 860–960 MHz family but reader power/channel rules are regional.

| Region | Band (approx.) | Power / access rule |
|---|---|---|
| US / FCC (Part 15) | 902–928 MHz | ≤1 W conducted, ≤4 W EIRP, frequency hopping |
| Europe / ETSI EN 302 208 | 865–868 MHz | ≤2 W ERP, **Listen-Before-Talk** |
| Japan | ~916–921 MHz (hist. 952–954) | Region-specific power/channel rules |
| China / others | vary (e.g. 920–925 MHz) | National regulator (MIIT, etc.) |

Practical consequence: design for the deployment region, or use globally-certified
multi-region reader hardware. Tags are generally broadband enough to work across the family;
the *reader* is the region-constrained component.

## Selecting a band (decision heuristics)
1. **Tap-to-authenticate, ID/access, payment, near the body, want short controlled range** →
   **HF (13.56 MHz)**, usually ISO/IEC 14443 (or 15693 for slightly longer vicinity reads).
2. **Bulk reading many items fast, supply chain / inventory / retail, metres of range** →
   **passive UHF**, ISO/IEC 18000-63 / EPC Gen2.
3. **Reading through liquids/tissue or around the body at very short range, animal ID,
   immobilizer** → **LF**.
4. **Real-time location, long range, onboard sensors, active beaconing** → **active/semi-passive
   UHF/microwave** (433 MHz or 2.45 GHz). (RTLS depth → asset-tracking sibling.)
5. **Items are wet or metallic** → expect range loss on UHF/microwave; budget for on-metal /
   liquid-tolerant tags, or prefer HF/LF for very short range.
