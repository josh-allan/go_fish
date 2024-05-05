```#graph TD;
    A[Browser] --> B{Cache};
    B -->|Found| C[Use Cached];
    B -->|Not Found| D[Resolve];
    D --> E[Local DNS Resolver];
    E -->|Found| F[Use Local DNS Cache];
    E -->|Not Found| G[Query Root DNS Servers];
    G --> H[Query TLD DNS Servers];
    H --> I[Query Authoritative DNS Servers];
    I --> J[Retrieve IP Address];
    J --> K[Return IP to Resolver];
    K --> L[Return IP to Browser];
    F --> L;
    L --> M[Use IP]; terraforming-mars
```
