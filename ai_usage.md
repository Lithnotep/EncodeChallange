1. First I had copilot make me a very simple GO skeleton for json/csv reading, with testing prepared to get started right away. as well as make sure project can be run when pulled down easily. Added the decode.json and encodes.csv as well.
2. Next I gave a description to copilot regarding the data and my plan. focusing on this idea "Use JSON streaming decoder instead of loading all 10k records into memory" 
3. Following that planning I moved to seperate the functionality to readers and aggregators. along with their test files. Most of my time spent double checking code manually and making sure I was moving in the right direction.(Claude 4 was very effective here. creating the packages and testing. INCLUDING a relatively complete understanding of the aggregation. Commiting this work but plan to change it to a more dynamic approach.)
4. Next steps will include Testing review as there appear to be unmapped bitly encodes.
# AI Usage Documentation

This document outlines how AI tools (GitHub Copilot and Claude) were used throughout the development of the EncodeChallenge project.

## Development Timeline

### 1. Initial Project Setup
**Tool Used:** GitHub Copilot  
**Objective:** Create a foundational Go project structure  

- Generated simple Go skeleton for JSON/CSV reading
- Set up testing framework for immediate TDD workflow
- Ensured project can be easily run when pulled down
- Added sample data files (`decodes.json` and `encodes.csv`)

### 2. Architecture Planning & Streaming Implementation
**Tool Used:** GitHub Copilot  
**Objective:** Design memory-efficient data processing  

**Key Focus:** *"Use JSON streaming decoder instead of loading all 10k records into memory"*

- Implemented streaming JSON processing for large datasets
- Avoided memory-intensive approaches that load entire files
- Designed for scalability with larger data volumes

### 3. Modular Code Organization  
**Tool Used:** Claude 4  
**Objective:** Separate concerns into testable modules

- Created separate `reader` and `aggregator` packages
- Developed comprehensive test suites for each module
- Implemented Test-Driven Development (TDD) approach
- **Result:** Claude effectively created packages with complete aggregation understanding

**Note:** Significant manual code review and validation was performed to ensure correct implementation direction.

### 4. Dynamic Command-Line Interface
**Tool Used:** Claude 4  
**Objective:** Add flexible year-based filtering  

**Prompt:** *"I want to be able to run a command in the terminal that will run the aggregator with dynamic settings. First option: I want to be able to send a specific year and have it get the clicks for that year. It should default to 2021 if no year is sent."*

**Implementation:**
- Added command-line flag support with Go's `flag` package
- Default year: 2021 (as per project requirements)
- Dynamic year filtering: any year can be queried
- Special case: `year=0` processes ALL years
- **Result:** Complete solution meeting project requirements

### 5. Final Output Formatting
**Tool Used:** Claude 4  
**Objective:** Create clean JSON-formatted output  

**Prompt:** *"I want the final log to be a super simple breakdown of the long URLs and the amount of clicks for the specified year. Right before this log, note that shortlinks without mapping are left out."*

**Example Format Requested:**
```json
[{"https://google.com": 3}, {"https://www.twitter.com": 2}]
```

**Implementation:**
- Clean JSON-like output format
- Automatic filtering of unmapped shortlinks
- Clear notification about excluded data
- Dynamic shortlink detection based on actual unmapped data

### 6. Data-Driven Shortlink Detection
**Tool Used:** Claude 4  
**Objective:** Replace hardcoded shortlink domains with dynamic detection

**Enhancement:** Instead of hardcoding shortlink domains (bit.ly, es.pn, etc.), the system now dynamically identifies shortlinks based on actual unmapped bitlinks encountered in the data.

**Benefits:**
- Adapts to new shortlink services automatically
- Works with any data source
- No maintenance of hardcoded domain lists

