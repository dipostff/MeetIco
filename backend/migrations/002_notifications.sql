CREATE TABLE notifications (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
  channel    TEXT NOT NULL,
  type       TEXT NOT NULL,
  send_at    TIMESTAMPTZ NOT NULL,
  sent       BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX notifications_pending_idx
  ON notifications(send_at, sent)
  WHERE sent = FALSE;
