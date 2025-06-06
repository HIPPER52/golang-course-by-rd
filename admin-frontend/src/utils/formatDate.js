export function formatDate(dateStr) {
    if (!dateStr) return 'N/A'
  
    const date = new Date(dateStr)
    return date.toLocaleString('en-GB', {
      year: 'numeric',
      month: 'short',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
