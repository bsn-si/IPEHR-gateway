![HMS gateway](https://user-images.githubusercontent.com/8058268/171113330-95e3a816-6b9c-4c09-83b1-c805f3feba48.png)

## REST API

At this point in the IPEHR gateway application, the minimum basic version of the REST API has been implemented to support document handling according to the latest stable version of the openEHR specification - 1.0.2  
<https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2>

In the next steps, the API will be supplemented with other methods to fully comply with the openEHR specifications, as well as additional ones, such as user and organization management and document access rights.

## Implementation

The code in which the REST API is implemented is located in the `src/pkg/api` directory
