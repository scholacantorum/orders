export default function (o) {
  const fd = new FormData()
  fd.append('source', 'inperson')
  fd.append('name', o.name)
  fd.append('email', o.email)
  o.lines.forEach((ol, idx) => {
    const prefix = `line${idx + 1}`
    fd.append(`${prefix}.product`, ol.product)
    fd.append(`${prefix}.quantity`, ol.quantity)
    fd.append(`${prefix}.price`, ol.price)
    if (ol.used) fd.append(`${prefix}.used`, ol.used)
    if (ol.usedAt) fd.append(`${prefix}.usedAt`, ol.usedAt)
  })
  fd.append('payment1.type', o.payments[0].type)
  if (o.payments[0].subtype) fd.append('payment1.subtype', o.payments[0].subtype)
  if (o.payments[0].method) fd.append('payment1.method', o.payments[0].method)
  fd.append('payment1.amount', o.payments[0].amount)
  return fd
}
