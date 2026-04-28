export {};

declare global {
  const Plotly: any;

  interface Window {
    APP_CONFIG: {
      apiBase: string;
    };
  }
}