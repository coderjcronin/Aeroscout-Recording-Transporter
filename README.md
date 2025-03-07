# Aeroscout Recording Transporter (ART)

## What is this?
This program, written in Go, is meant to be a utiliy for preparing recording sessions made by the Aeroscout Location Engine (ALE) to minimize their footprint for transport as well as restore their layout for recording analysis/reporting.

This is accomplish by iterating through the recording folders within the session, deleting all but the first instance of the RadioMaps and Maps folders (since these folders are identical through recordings), then zip archiving the recording session when ready.

## Why Go?
Because I know maybe four languages well enough to do this:
 - Python - Requires Python to be installed and that's going to be a headache for most people who need this on the fly on client servers
 - Powershell - Requires code signing or the client being really cool with running unsigned code (hint: they won't be)
 - C++ - Just no. No. No no no no no.
 - Go - Well... winner winner.

 ## Why do this?
 ALE recording sessions are composed of folders containing each distinctive one-click recording. Each folder then contains copies of the ALE's maps and radio maps. On smaller ALE installations this is not an issue, however, on larger enterprise or departmental systems this means copying large amounts of identical data multiple times.

 While this doesn't sound like a lot, here's a real world example:
Client A is an enterprise client whose ALE is tracking 5 facilities with an average of 2 units per facility (10 tracked areas). They are adding a new unit to an existing facility, requiring recordings of the base floor of the facility plus two tower floors. The base floor takes 100 recordings while the two tower floors are 40 recordings each. The total size of the three sessions will be the sum of Maps, RadioMaps, and a recording multiplied 180 times. In this example, RadioMaps and Maps weigh in at 600 MB while an average recording is 37.3 MB; this bring the total to 112 GB ((180 * (600 + 37.3)) / 1024). 
 
 Taking an large chunk of the client's data drive alone could pose a problem (causing a low disk space alert or requiring them to allocate additional space which will incur cost), however, these recordings also need to be uploaded to storage for backup and retrieved as needed for analysis causing cost to be incurred for bandwidth, download time, and storage.

 Using ART in the above Client A example, we get rid of 177 of the 180 copies of the Radio Maps and Maps, reducing size from 112 GB to 8.3 GB (180 times the average size of recording plus three copies of the maps) for a size reduction of 93%. This will ease file transfer and cloud storage as well as mean less downtime for technicians working with the recordings waiting for file operations to be executed.

Testing has found that the ALE only care about RadioMaps and Maps in the first recording folder (as of 5.7, at least). So we will concentrate on cleaning up the recording session folder and potentially zip archiving the session for transfer.
 ## Why such verbose comments and ReadMe?
 This application is meant to be used on healthcare client machines. By showing, verbosely, what the application is and what it does I can create accountability and understanding for it's operation. This client has no file transport ability outside the host system, no remote administrative toolkit (RAT), nor does it inject itself or other code into any other part of the client system. By being able to review the code as written prior to compilation and verify that operation my hope is to instill confidence that the best interest of the client's IT security and safety is shown.

### Assumptions
1. This will only be run within an ALE recording session folder
2. ~This will be run only twice (once to prep a recording prior to "transport", once to restore the recording for analysis)~
3. This utility should be "frozen" so a copy can be left within the recording session folder with a small footprint (kb, not mb)
4. CLI is fine with minimal user interface (I really don't feel like doing flags for true CLI interface but I will if there's a need)

### Recording Session Tree
Folders for recording sessions are typically laid out as such:
- 2000-01-01 Client Facility Area
  - Datetime stamp
    - Devices
    - Maps
    - RadioMaps
    - MAPPlaybackFile.dat
    - _Other files that we don't need to worry about_
  - Datatime Stamp
    - Maps
    - RadioMaps
    - MAPPlaybackFile.dat
    - _Other files that we still don't worry about_
  - _Repeat for other recordings if present_

 ### TODO
 - Refactor as testing showed the ALE only cares about Maps and RadioMaps in the **first recording folder**
   - Remove Processing for Analysis
   - Rename "Processing for Transport" to "Clean Recording Session"
   - Add "Archive for Transport" (testing found zip archive further reduced size to < 10%, ideal for easy upload and storage)