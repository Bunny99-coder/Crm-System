// File: src/util/auth.ts
export const getUserIdFromToken = (token: string | null): number | null => {
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.user_id || null;
  } catch (e) {
    console.error("Failed to parse token:", e);
    return null;
  }
};


export const getRoleFromToken = (token: string | null): number | null => {
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.role_id || null;
  } catch (e) {
    return null;
  }
};