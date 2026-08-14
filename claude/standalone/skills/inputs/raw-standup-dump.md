notes week of 6/29

mon standup - kavya: search latency spike friday night, p99 800ms. rolled back the analyzer change
TODO file incident writeup
carlos out til thursday
- pricing meeting moved to weds
decided were going with the usage-based tier for the API, flat tier killed
todo: update pricing page draft (sam has it)
what if we exposed the latency dashboard to customers? maybe as a trust page thing
kavya says rollback fixed p99, back to 120ms
weds - pricing mtg notes: legal wants 30 day notice for price changes, EU customers 60?? check with dana
usage based tier CONFIRMED going in the sept release
TODO file incident report for fridays latency thing
sam: pricing page draft blocked on final tier names
tier names: builder / scale / enterprise (marketing hasnt signed off)
idea - trust page could also show uptime, status.io does this
thurs: carlos back, says the analyzer change actually shipped to EU cluster too, never rolled back there
!!! EU p99 still degraded, carlos filing ticket
dana confirmed: EU price change notice is 60 days, its in the DPA addendum
decided: hold pricing page until marketing signs off tier names
q: does the sept release date slip if usage metering isnt done?
metering work estimate: 3 wks (jt), release is 5 wks out so ok
TODO ping marketing re tier names (sam owns)
