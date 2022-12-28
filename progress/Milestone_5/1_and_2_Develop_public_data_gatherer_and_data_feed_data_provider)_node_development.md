We developed a service that allows making statistical data available to be collected and processed by the IPEHR system. The service implements an open API with specified metrics. The data is collected and processed by accessing [IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes) smart contracts.
  
Medical statistics can be collected by the service through:

  

-   transaction analysis of contracts [IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes)
    
-   periodic direct invocation of contract methods [IPEHR-blockchain-indexes](https://github.com/bsn-si/IPEHR-blockchain-indexes)
    
-   making AQL queries to IPEHR-gateway
    

  

For demonstration purposes, the following metrics are implemented:

  

-   number of patients registered in the system overall time
    
-   number of patients logged in the system for a specified month
    
-   number of EHR documents registered in the system for all time
    
-   number of EHR documents registered in the system for a given month
    

  

Please note that the following metrics are samples that are made to demonstrate the possibilities of MVP. At this point, we are providing access only to statistical data. With the future development of the project when it will be monetized with its own tokens, we will be able to provide access to a much wider range of data and contracts.

  

A [link](https://stat.ipehr.org/swagger/index.html) to API documentation.

A [link](https://github.com/bsn-si/IPEHR-stat) for more information and manuals.