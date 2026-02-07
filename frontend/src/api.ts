import { ref, computed, onMounted } from 'vue'

const API_BASE = ''

// === 接口类型 ===

export interface FormatInfo {
  id: string
  name: string
  kind: number // 1=拼音 2=五笔 3=词组
  ext: string  // 支持的文件扩展名
  canImport: boolean
  canExport: boolean
}

export interface UploadResult {
  id: string
  filename: string
  size: number
}

export interface EncoderConfig {
  type: string        // "pinyin" | "wubi" | "none"
  schema: string      // "86" | "98" | "06"
  codeTableFileId: string
  useAABC: boolean
}

export interface FilterConfig {
  minLength: number
  maxLength: number
  minFrequency: number
  filterEnglish: boolean
  filterNumber: boolean
  customRules: string[]
}

export interface CustomFieldConfig {
  type: string
  pinyinSeparator: string
  pinyinPrefix: string
  pinyinSuffix: string
  literal: string
}

export interface CustomFormatConfig {
  kind: string
  encoding: string
  fields: CustomFieldConfig[]
  sortByCode: boolean
  commentPrefix: string
  startMarker: string
}

export interface ConvertRequest {
  fileId: string
  inputFormat: string
  outputFormat: string
  inputCustom?: CustomFormatConfig
  outputCustom?: CustomFormatConfig
  encoder?: EncoderConfig
  filter?: FilterConfig
}

export interface ConvertResult {
  outputPath: string
  stats: {
    inputEntries: number
    outputEntries: number
    filteredOut: number
  }
}

// === API 调用 ===

export async function fetchFormats(): Promise<FormatInfo[]> {
  const res = await fetch(`${API_BASE}/api/formats`)
  return res.json()
}

export async function uploadFile(file: File): Promise<UploadResult> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await fetch(`${API_BASE}/api/upload`, {
    method: 'POST',
    body: formData,
  })
  return res.json()
}

export async function convertFile(req: ConvertRequest): Promise<ConvertResult> {
  const res = await fetch(`${API_BASE}/api/convert`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || '转换失败')
  }
  return res.json()
}

// === 工具函数 ===

export function formatKindLabel(kind: number): string {
  switch (kind) {
    case 1: return '拼音'
    case 2: return '五笔'
    case 3: return '词组'
    default: return '未知'
  }
}

/** 根据文件扩展名自动匹配输入格式 */
export function matchFormatByExt(filename: string, formats: FormatInfo[]): string | null {
  const lower = filename.toLowerCase()
  // ext 包含前导点，如 ".scel"、".dict.yaml"，用 endsWith 匹配
  const match = formats.find(f => f.canImport && f.ext && lower.endsWith(f.ext.toLowerCase()))
  return match ? match.id : null
}

// === 全局状态 composable ===

export function useFormats() {
  const formats = ref<FormatInfo[]>([])
  const loading = ref(false)

  const importFormats = computed(() =>
    formats.value.filter(f => f.canImport)
  )
  const exportFormats = computed(() =>
    formats.value.filter(f => f.canExport)
  )

  async function load() {
    loading.value = true
    try {
      formats.value = await fetchFormats()
    } finally {
      loading.value = false
    }
  }

  onMounted(load)

  return { formats, importFormats, exportFormats, loading, reload: load }
}
