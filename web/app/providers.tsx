'use client';

import {NextUIProvider, Link} from '@nextui-org/react';
import {useRouter} from 'next/navigation'

export function Providers({children}: { children: React.ReactNode }) {
  const router = useRouter();
  return (
    <NextUIProvider 
      navigate={router.push}
      // navigate={function push(path: string): void {
      //   if (path === undefined) {
      //     // throw new Error('missing path!!!!');
      //   }
      //   console.log(path)
      //   router.push(path)
      // }}
    >
      {children}
    </NextUIProvider>
  )
}