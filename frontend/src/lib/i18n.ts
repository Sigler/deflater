// The active language. Future languages: add strings/<lang>.ts with the
// same shape, then pick here (for example from navigator.language).
import { en } from "./strings/en";

export const S = en;
export type Strings = typeof en;
