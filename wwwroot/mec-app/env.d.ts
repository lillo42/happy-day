/// <reference types="vite/client" />
export interface ImportMetaEnv {
  readonly VITE_MEC_API: string
}

export interface ImportMeta {
  readonly env: ImportMetaEnv
}
