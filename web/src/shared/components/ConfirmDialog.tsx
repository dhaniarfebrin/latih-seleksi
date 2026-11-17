import { useEffect, useRef, type ReactNode } from 'react'

interface Props {
  open: boolean
  title: string
  children: ReactNode
  confirmLabel: string
  busy?: boolean
  error?: string
  onConfirm: () => void
  onCancel: () => void
}

// Native <dialog>: showModal() traps focus, Esc fires `cancel`, ::backdrop dims the page.
export function ConfirmDialog({ open, title, children, confirmLabel, busy, error, onConfirm, onCancel }: Props) {
  const ref = useRef<HTMLDialogElement>(null)
  const cancelRef = useRef<HTMLButtonElement>(null)
  useEffect(() => {
    const d = ref.current
    if (!d) return
    if (open && !d.open) {
      d.showModal()
      cancelRef.current?.focus()
    } else if (!open && d.open) d.close()
  }, [open])
  return (
    <dialog
      ref={ref}
      className="dlg"
      aria-labelledby="dlg-title"
      onCancel={(e) => {
        e.preventDefault()
        if (!busy) onCancel()
      }}
    >
      <div className="body">
        <h2 id="dlg-title">{title}</h2>
        {children}
        {error && (
          <p role="alert" className="err">
            {error}
          </p>
        )}
        <div className="row">
          <button ref={cancelRef} type="button" className="btn ghost" disabled={busy} onClick={onCancel}>
            Batal
          </button>
          <button type="button" className="btn danger" disabled={busy} onClick={onConfirm}>
            {confirmLabel}
          </button>
        </div>
      </div>
    </dialog>
  )
}
