import { atom } from "recoil";

export const authState = atom({
  key: 'Auth',
  default: {
    isAuthenticated: false
  }
});
