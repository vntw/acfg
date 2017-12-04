export const SESSION_BOOKING = 0;
export const SESSION_PRACTICE = 1;
export const SESSION_QUALIFYING = 2;
export const SESSION_RACE = 3;

export function getSessionName(session) {
  switch (session) {
    case SESSION_BOOKING:
      return 'Booking';
    case SESSION_PRACTICE:
      return 'Practice';
    case SESSION_QUALIFYING:
      return 'Qualifying';
    case SESSION_RACE:
      return 'Race';
    default:
      return 'Unknown';
  }
}

export function getShortSessionName(session) {
  return getSessionName(session)[0];
}
