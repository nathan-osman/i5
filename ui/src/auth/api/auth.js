import { atom } from "recoil";

const authAtom = atom({
  key: 'Auth',
  default: {
    isAuthenticated: false
  }
});

export { authAtom };
