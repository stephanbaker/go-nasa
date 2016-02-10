# go-nasa

The purpose of this repository is to host a series of packages aimed at assisting in the retrieval of data from NASA's public APIs using Google's Go language.  In order to use these APIs, you must register for an API key at https://api.nasa.gov/index.html#apply-for-an-api-key.

# Package: apod

Assists in the retrieval of data from NASA's Astronomy Picture of the Day (APOD) API.  See https://api.nasa.gov/api.html#apod for more information.

#### GetEntry
GetEntry queries the NASA APOD API and returns the result for the given parameters.  This function expects the api url, date of the entry, whether or not to provide the high-defenition URL for entry, and the API key to be used for the request.

#### GetEntryForDate
GetEntryForDate queries the NASA APOD API and returns the result for the given parameters.  This function expects the date of the entry to be provided.

#### GetEntryForToday
GetEntryForToday queries the NASA APOD API and returns the today's entry.
