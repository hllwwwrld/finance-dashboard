// DELETE /api/payments/[id] - удалить платеж

import { NextRequest, NextResponse } from 'next/server'

export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params

    if (!id) {
      return NextResponse.json(
        { error: 'ID платежа не указан' },
        { status: 400 }
      )
    }

    // TODO: Замените на удаление из реальной базы данных
    // await db.payments.delete({ where: { id } })

    return NextResponse.json({ success: true })
  } catch (error) {
    console.error('Ошибка при удалении платежа:', error)
    return NextResponse.json(
      { error: 'Не удалось удалить платеж' },
      { status: 500 }
    )
  }
}
