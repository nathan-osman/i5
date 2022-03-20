import { atom } from "recoil";

const authState = atom({
  key: 'Auth',
  default: {
    isAuthenticated: false
  }
});
