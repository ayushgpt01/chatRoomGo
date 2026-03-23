interface DevImportMetaEnv {
	readonly VITE_API_URL: string;
	readonly VITE_ENV: "development";
	readonly VITE_WS_URL: string;
}

interface ProdImportMetaEnv {
	readonly VITE_ENV: "production";
	readonly VITE_API_ROUTE: string;
}

interface ImportMeta {
	readonly env: DevImportMetaEnv | ProdImportMetaEnv;
}
