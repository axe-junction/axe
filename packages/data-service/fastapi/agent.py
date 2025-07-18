from agno.agent import Agent
from agno.models.google import Gemini
from agno.utils.media import download_file
from agno.media import File
from pathlib import Path
from pydantic import BaseModel, Field
import asyncio
import requests
import json
import uuid
from code.app.models.station import placeInDb

def geocode_station(station_name):
    """Geocode a station using OpenStreetMap Nominatim API"""
    try:
        # Format the query to include Algeria for better results
        query = f"{station_name} Algeria"
        url = f"https://nominatim.openstreetmap.org/search?q={query}&format=json&limit=1"
        
        # Make the API request
        response = requests.get(url, headers={'User-Agent': 'SNTF-Data-Extractor'})
        response.raise_for_status()
        
        data = response.json()
        
        if data:
            result = data[0]
            return {
                'latitude': float(result['lat']),
                'longitude': float(result['lon']),
                'display_name': result['display_name'],
                'confidence': float(result.get('importance', 0.5))
            }
        else:
            print(f"No geocoding results found for: {station_name}")
            return None
            
    except Exception as e:
        print(f"Geocoding failed for {station_name}: {e}")
        return None

# Station Extraction Agent using Gemini 2.0 Flash
station_extraction_agent = Agent(
    name="SNTF Station Extraction Agent",
    role="Extract station names from SNTF schedule images using Gemini 2.0 Flash",
    model=Gemini(
        id="gemini-2.0-flash", 
        api_key="AIzaSyCtGGY1Sb6RVSEiGuBIuXjEbDlcCrQXgwU",
    ),
    description="""
    SNTF Station Extraction Agent powered by Gemini 2.0 Flash.
    
    I specialize in:
    • Analyzing SNTF schedule images using vision capabilities
    • Extracting ALL station names from railway schedule tables
    • Identifying station sequences and route information
    • Handling French and Arabic text in schedule documents
    • Providing clean, structured station name lists
    
    I process ALL three SNTF railway lines:
    - Est Line (Eastern Algeria): All stations from Algiers to Thenia/Annaba
    - Ouest Line (Western Algeria): All stations from Algiers to El Affroun/Oran
    - Zeralda Line (Algiers suburban): All stations from Algiers to Zeralda
    """,
    instructions="""
    You are a station name extraction specialist using Gemini 2.0 Flash vision capabilities.

    ## PRIMARY OBJECTIVE:
    Extract ALL station names from SNTF schedule images and return them as clean, structured lists.

    ## EXTRACTION PROCESS:

    ### 1. IMAGE ANALYSIS
    - Analyze the provided SNTF schedule images using vision understanding
    - Identify table structures containing station information
    - Extract text from all visible schedule elements
    - Focus specifically on station names, not schedule times

    ### 2. STATION NAME EXTRACTION
    - Look for station names in BOTH French and Arabic text
    - Extract ALL station names from each railway line
    - Identify station sequences and route order
    - Handle variations in station name formatting
    - Clean and normalize station names

    ### 3. OUTPUT FORMAT
    Return station names in this JSON structure:
    {
        "Est": ["Alger", "Station1", "Station2", ..., "Thenia"],
        "Ouest": ["Alger", "Station1", "Station2", ..., "El Affroun"],
        "Zeralda": ["Alger", "Station1", "Station2", ..., "Zeralda"]
    }

    ## CRITICAL REQUIREMENTS:
    - Extract EVERY visible station name from all three images
    - Maintain station order along each route
    - Handle multilingual content (French/Arabic)
    - Provide clean, consistent station names
    - Include intermediate stops, not just terminus stations

    ## QUALITY STANDARDS:
    - Extract at least 10-15 stations per line
    - Ensure completeness - miss no visible stations
    - Maintain accuracy in station name spelling
    - Handle abbreviated or partial station names

    Focus on COMPLETENESS. Your job is to extract station names only - geocoding will be handled separately using OpenStreetMap API.
    """,
    markdown=True,
)

# Unified SNTF Data Processing Agent
unified_sntf_agent = Agent(
    name="SNTF Data Processing Agent",
    role="Process extracted stations and create complete railway data with geocoding",
    model=Gemini(
        id="gemini-2.0-flash", 
        api_key="AIzaSyCtGGY1Sb6RVSEiGuBIuXjEbDlcCrQXgwU",
    ),
    description="""
    SNTF Data Processing Agent that takes extracted station names and creates complete railway data.
    
    I handle:
    • Processing lists of extracted station names
    • Coordinating geocoding using OpenStreetMap Nominatim API
    • Structuring data according to placeInDb model
    • Validating and cleaning processed data
    • Ensuring comprehensive coverage of the SNTF network
    """,
    instructions="""
    You are a data processing specialist that takes extracted station names and creates complete railway data.

    ## WORKFLOW:

    ### 1. PROCESS STATION LISTS
    - Take the extracted station names from the extraction agent
    - Validate and clean station names
    - Organize by railway line (Est, Ouest, Zeralda)

    ### 2. DATA STRUCTURING
    - Structure data according to placeInDb model with these EXACT requirements:
    - id: Generate valid UUID strings using uuid4() format
    - name: Station name as extracted
    - coordinates: {"latitude": float, "longitude": float}
    - type: Always "railway_station"
    - region: Always "Algeria"

    ### 3. IMPORTANT UUID REQUIREMENTS:
    - Generate proper UUIDs in format: "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
    - Use valid UUID4 format only
    - Do NOT use placeholder strings like "placeholder_station_123"
    - Each station must have a unique, valid UUID

    ### 4. QUALITY ASSURANCE
    - Ensure all stations have valid coordinates
    - Check for duplicates and missing data
    - Validate geographic locations within Algeria
    - Provide confidence scores and data quality metrics

    ## OUTPUT REQUIREMENTS:
    Return complete railway station data in JSON format matching the placeInDb model structure with valid UUIDs.
    """,
    response_model=placeInDb,
    markdown=True,
)

def download_sntf_images():
    """Download SNTF schedule images locally"""
    base_path = Path(__file__).parent.joinpath("sntf_images")
    base_path.mkdir(exist_ok=True)
    
    urls = {
        "Est": "https://www.sntf.dz/documents/Est.jpg",
        "Ouest": "https://www.sntf.dz/documents/Ouest.jpg", 
        "Zeralda": "https://www.sntf.dz/documents/Zeralda.jpg"
    }
    
    downloaded_files = []
    
    for line_name, url in urls.items():
        try:
            file_path = base_path.joinpath(f"{line_name}.jpg")
            download_file(url, str(file_path))
            downloaded_files.append(File(filepath=file_path))
            print(f"Downloaded {line_name} line image to {file_path}")
        except Exception as e:
            print(f"Failed to download {line_name} line image: {e}")
    
    return downloaded_files

def save_stations_to_file(stations_data, filename="sntf_stations.json"):
    """Save extracted station data to JSON file"""
    try:
        output_path = Path(__file__).parent.joinpath(filename)
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(stations_data, f, indent=2, ensure_ascii=False)
        print(f"Station data saved to: {output_path}")
        return output_path
    except Exception as e:
        print(f"Failed to save station data: {e}")
        return None

def structure_station_data(geocoded_stations):
    """Structure station data according to placeInDb model"""
    structured_data = []
    
    for line_name, stations in geocoded_stations.items():
        for station in stations:
            # Create proper placeInDb structure
            station_data = {
                "id": station['id'],  # Already has proper UUID
                "name": station['name'],
                "coordinates": station['coordinates'],
                "type": station['type'],
                "region": station['region']
            }
            structured_data.append(station_data)
    
    return structured_data

def clean_json_response(response_content):
    """Clean JSON response by removing markdown code blocks"""
    try:
        # Remove markdown code block markers
        content = response_content.strip()
        
        # Remove ```json at the beginning
        if content.startswith('```json'):
            content = content[7:]  # Remove ```json
        elif content.startswith('```'):
            content = content[3:]   # Remove ```
        
        # Remove ``` at the end
        if content.endswith('```'):
            content = content[:-3]
        
        # Strip whitespace
        content = content.strip()
        
        return content
    except Exception as e:
        print(f"Error cleaning JSON response: {e}")
        return response_content

if __name__ == "__main__":
    try:
        # Download images first
        print("Downloading SNTF schedule images...")
        downloaded_files = download_sntf_images()
        
        if downloaded_files:
            # Step 1: Extract station names using Gemini 2.0 Flash
            print("Extracting station names from images...")
            station_extraction_response = station_extraction_agent.run(
                """Extract ALL station names from the provided SNTF schedule images.
                
                Analyze all three railway lines:
                1. Est Line (Eastern Algeria) - Extract ALL stations from Algiers to Thenia and beyond
                2. Ouest Line (Western Algeria) - Extract ALL stations from Algiers to El Affroun and beyond
                3. Zeralda Line (Algiers suburban) - Extract ALL stations from Algiers to Zeralda
                
                Return the station names in JSON format:
                {
                    "Est": ["station1", "station2", ...],
                    "Ouest": ["station1", "station2", ...],
                    "Zeralda": ["station1", "station2", ...]
                }
                
                Focus on extracting EVERY visible station name, not just major terminals.""",
                files=downloaded_files
            )
            
            print("Station extraction completed.")
            
            # LOG: Print raw extraction response
            print("\n" + "="*50)
            print("RAW EXTRACTION RESPONSE:")
            print("="*50)
            print(f"Response type: {type(station_extraction_response)}")
            print(f"Response content: {station_extraction_response.content}")
            print(f"Response length: {len(station_extraction_response.content)}")
            print("="*50 + "\n")
            
            # Step 2: Clean and parse extracted station names
            try:
                # Clean the response to remove markdown code blocks
                cleaned_content = clean_json_response(station_extraction_response.content)
                
                print("\n" + "-"*30)
                print("CLEANED JSON CONTENT:")
                print("-"*30)
                print(f"Cleaned content: {cleaned_content[:500]}...")
                print("-"*30 + "\n")
                
                # Parse the cleaned JSON
                station_data = json.loads(cleaned_content)
                print(f"✓ Successfully parsed JSON")
                print(f"Extracted stations: {station_data}")
                
                # LOG: Detailed analysis of extracted data
                print("\n" + "-"*30)
                print("EXTRACTED DATA ANALYSIS:")
                print("-"*30)
                for line_name, stations in station_data.items():
                    print(f"{line_name} line: {len(stations)} stations")
                    if stations:
                        print(f"  First station: {stations[0]}")
                        print(f"  Last station: {stations[-1]}")
                        print(f"  All stations: {stations}")
                    else:
                        print(f"  WARNING: No stations found for {line_name} line!")
                print("-"*30 + "\n")
                
            except json.JSONDecodeError as e:
                print(f"✗ Failed to parse station extraction response as JSON: {e}")
                print(f"Raw content: {station_extraction_response.content[:500]}...")
                
                # Try to extract data manually if JSON parsing fails
                content = station_extraction_response.content
                print("\n" + "-"*30)
                print("MANUAL EXTRACTION ATTEMPT:")
                print("-"*30)
                
                # Look for common patterns in the response
                if "Est" in content or "Ouest" in content or "Zeralda" in content:
                    print("✓ Found railway line names in response")
                else:
                    print("✗ No railway line names found in response")
                    
                if "Alger" in content or "Algiers" in content:
                    print("✓ Found Algiers references in response")
                else:
                    print("✗ No Algiers references found in response")
                    
                print("-"*30 + "\n")
                station_data = {}
            
            except Exception as e:
                print(f"✗ Unexpected error parsing response: {e}")
                station_data = {}
            
            # LOG: Check if we have any data to process
            if not station_data:
                print("⚠️  WARNING: No station data extracted. Cannot proceed with geocoding.")
                print("This could be due to:")
                print("1. Images not downloading properly")
                print("2. Vision model not recognizing text in images")
                print("3. Images not containing expected schedule data")
                print("4. JSON parsing issues")
                
                # Check downloaded files
                print("\n" + "-"*30)
                print("DOWNLOADED FILES CHECK:")
                print("-"*30)
                for file in downloaded_files:
                    print(f"File: {file.filepath}")
                    if file.filepath.exists():
                        file_size = file.filepath.stat().st_size
                        print(f"  Size: {file_size} bytes")
                        print(f"  Exists: ✓")
                    else:
                        print(f"  Exists: ✗")
                print("-"*30 + "\n")
                
                exit()
            
            # Step 3: Geocode stations using OpenStreetMap API
            print("Geocoding stations using OpenStreetMap API...")
            geocoded_stations = {}
            total_stations = 0
            
            for line_name, stations in station_data.items():
                print(f"\nProcessing {line_name} line ({len(stations)} stations)...")
                geocoded_stations[line_name] = []
                
                for i, station in enumerate(stations, 1):
                    print(f"  [{i}/{len(stations)}] Geocoding: {station}")
                    geo_result = geocode_station(station)
                    
                    if geo_result:
                        # Generate a proper UUID for each station
                        station_uuid = str(uuid.uuid4())
                        geocoded_stations[line_name].append({
                            'id': station_uuid,
                            'name': station,
                            'coordinates': {
                                'latitude': geo_result['latitude'],
                                'longitude': geo_result['longitude']
                            },
                            'type': 'train',
                            'region': 'Algeria',
                            'display_name': geo_result['display_name'],
                            'confidence': geo_result['confidence'],
                            'line': line_name
                        })
                        print(f"    ✓ {station}: {geo_result['latitude']}, {geo_result['longitude']}")
                        total_stations += 1
                    else:
                        print(f"    ✗ Failed to geocode: {station}")
                
                print(f"  Completed {line_name} line: {len(geocoded_stations[line_name])} stations geocoded")
            
            print(f"\nTotal stations geocoded: {total_stations}")

            # Step 4: Structure final data directly (without problematic agent)
            print("Structuring final data...")
            structured_stations = structure_station_data(geocoded_stations)
            
            # Step 5: Save to file
            print("Saving station data to file...")
            
            # Save the complete data with metadata
            final_data = {
                "metadata": {
                    "source": "SNTF Schedule Images",
                    "extraction_date": str(uuid.uuid4()),  # Use timestamp if needed
                    "total_stations": len(structured_stations),
                    "lines": list(geocoded_stations.keys())
                },
                "stations": structured_stations,
                "raw_geocoded_data": geocoded_stations
            }
            
            # Save to JSON file
            saved_file = save_stations_to_file(final_data, "sntf_stations_complete.json")
            
            # Also save just the structured stations for placeInDb compatibility
            placeInDb_compatible = {
                "stations": structured_stations
            }
            save_stations_to_file(placeInDb_compatible, "sntf_stations_placeInDb.json")
            
            print("Data processing completed successfully!")
            print(f"Total stations processed: {len(structured_stations)}")
            print(f"Lines processed: {list(geocoded_stations.keys())}")
            
            # Print sample data for verification
            if structured_stations:
                print("\nSample station data:")
                sample_station = structured_stations[0]
                print(f"ID: {sample_station['id']}")
                print(f"Name: {sample_station['name']}")
                print(f"Coordinates: {sample_station['coordinates']}")
                print(f"Type: {sample_station['type']}")
                print(f"Region: {sample_station['region']}")
            
        else:
            print("No images were successfully downloaded.")
            
    except Exception as e:
        print(f"An error occurred: {e}")
        import traceback
        traceback.print_exc()


