{
  fetchurl,
  version ? "2.8.1",
}:
fetchurl {
  url = "https://github.com/pfrest/pfSense-pkg-RESTAPI/releases/download/v${version}/openapi.json";
  hash = "sha256-Va7wPj+AsrYRGBz/ZEQChJKo10oFSkxAIFNDBAHcIOI=";
}
