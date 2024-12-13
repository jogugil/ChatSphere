/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_API_URL: string;
    readonly VITE_IP_SERVER_GOCHAT: string;
    readonly VITE_PORT_SERVER_GOCHAT: string;
    // Puedes agregar más variables de entorno según sea necesario
  }
  
  interface ImportMeta {
    readonly env: ImportMetaEnv;
  }