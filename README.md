# Aeroscout Recording Transporter (ART)

## What is this?
This program, written in Go, is meant to be a utiliy for preparing recording sessions made by the Aeroscout Location Engine (ALE) to minimize their footprint for transport as well as restore their layout for recording analysis/reporting.

This is accomplish by iterating through the recording folders within the session, deleting all but the first instance of the RadioMaps and Maps folders (since these folders are identical through recordings), then restoring the directory layout via symbolic links when needed for analysis.

## Why Go?
Because I know maybe four languages well enough to do this:
 - Python - Requires Python to be installed and that's going to be a headache for most people who need this on the fly on client servers
 - Powershell - Requires code signing or the client being really cool with running unsigned code (hint: they won't be)
 - C++ - Just no. No. No no no no no.
 - Go - Well... winner winner.

 ## Why do this?
 ALE recording sessions are composed of folders containing each distinctive one-click recording. Each folder then contains copies of the ALE's maps and radio maps. On smaller ALE installations this is not an issue, however, on larger enterprise or departmental systems this means copying large amounts of identical data multiple times.

 While this doesn't sound like a lot, here's a real world example:
 > Client A is running a VM instance of the ALE on Windows Server, with the ALE living on a separate data drive (E:) allocated to 100 GB. 
 > Client A is an enterprise client, running 4 facilities each containing 2 units with a moderate footprint in each unit. 
 > Radio Maps for Client A's sites total at approximately 500 MB. 
 > Client A is adding a site with a single large footprint, and requires new recording of this area for coverage and accuracy analysis. 
 > 100 recordings are made of the new unit footprint to ensure adequate coverage. Each recording, on it's own, is approximately 35 MB plus the copy of the ALE's maps and radio maps. 
 > The total session folder size will be 100 times the sum of the average recording size and the size of the ALE's radio maps/maps. This will be about 52.24 GB total, or over half the allocated drive size. 
 
 Taking an large chunk of the client's data drive alone could pose a problem (causing a low disk space alert or requiring them to allocate additional space which will incur cost), however, these recordings also need to be uploaded to storage for backup and retrieved as needed for analysis causing cost to be incurred for bandwidth, download time, and storage.

 Using ART in the above Client A example, we get rid of 99 of the 100 copies of the Radio Maps and Maps, reducing size from 52.24 GB to 3.9 GB (99 times the average size of recording plus the sum of the recording and a single copy of the maps) for a size reduction of 93%. This will ease file transfer and cloud storage as well as mean less downtime for technicians working with the recordings waiting for file operations to be executed.

 Since the Radio Maps and Maps folders are identiical across all recordings of a session (unless someone is editing maps in the ALE during recording, which is against Best Practices) we can simply iterate symbolic links (TODO: Make sure symbolic works... may have to switch to hard) in the other folders.

 ## Why such verbose comments and ReadMe?
 This application is meant to be used on healthcare client machines. By showing, verbosely, what the application is and what it does I can create accountability and understanding for it's operation. This client has no file transport ability outside the host system, no remote administrative toolkit (RAT), nor does it inject itself or other code into any other part of the client system. By being able to review the code as written prior to compilation and verify that operation my hope is to instill confidence that the best interest of the client's IT security and safety is shown.

### Assumptions
1. This will only be run within an ALE recording session folder
2. This will be run only twice (once to prep a recording prior to "transport", once to restore the recording for analysis)
3. This utility should be "frozen" so a copy can be left within the recording session folder with a small footprint (kb, not mb)
4. CLI is fine with minimal user interface (I really don't feel like doing flags for true CLI interface but I will if there's a need)

 ### TODO
 - ~~Complete initial prep functions (recording and analysis)~~
 - ~~Sanity check to prevent running this outside a recording session~~
 - Test to be sure symlinks are allowed for ALE analysis (may need hard links, not tested)
 - ~~Output log to a text file with time stamp for tracking~~