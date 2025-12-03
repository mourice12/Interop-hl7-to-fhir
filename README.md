# HL7 v2.x to FHIR R4 Converter

A Go-based converter that transforms HL7 v2.x messages into FHIR R4 Bundle resources.

## Features

- Parses HL7 v2.x messages (ADT, etc.)
- Converts to FHIR R4 resources:
  - Patient (from PID)
  - Encounter (from PV1)
  - Condition (from DG1)
  - AllergyIntolerance (from AL1)
  - Observation (from OBX)
- REST API endpoint
- Docker support

## Usage

### REST API

```bash
go run ./cmd/server
Then POST HL7 messages to http://localhost:8000/convert
Docker
docker build -t hl7-to-fhir .
docker run -p 8000:8000 hl7-to-fhir
License
MIT License - see LICENSE file