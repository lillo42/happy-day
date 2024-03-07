/* eslint-disable-next-line arrow-body-style */
export const spyGetter = <T, K extends keyof T>(target: jasmine.SpyObj<T>, key: K): jasmine.Spy => {
  return Object.getOwnPropertyDescriptor(target, key)?.get as jasmine.Spy;
};

/* eslint-disable-next-line arrow-body-style */
export const spySetter = <T, K extends keyof T>(target: jasmine.SpyObj<T>, key: K): jasmine.Spy => {
  return Object.getOwnPropertyDescriptor(target, key)?.set as jasmine.Spy;
};
