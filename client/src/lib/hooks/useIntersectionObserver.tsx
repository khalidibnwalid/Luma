// credits for https://usehooks-ts.com/react-hook/use-intersection-observer
// with some modifications to expose the ref and setRef

import { useEffect, useRef, useState } from 'react'

type State = {
  isIntersecting: boolean
  entry?: IntersectionObserverEntry
}

type UseIntersectionObserverOptions = {
  root?: Element | Document | null
  rootMargin?: string
  threshold?: number | number[]
  freezeOnceVisible?: boolean
  onChange?: (isIntersecting: boolean, entry: IntersectionObserverEntry) => void
  initialIsIntersecting?: boolean
}

type IntersectionReturn = [
  (node?: Element | null) => void,
  Element | null,
  boolean,
  IntersectionObserverEntry | undefined,
] & {
  setRef: (node?: Element | null) => void,
  ref: Element | null,
  isIntersecting: boolean
  entry?: IntersectionObserverEntry
}

export function useIntersectionObserver({
  threshold = 0,
  root = null,
  rootMargin = '0%',
  freezeOnceVisible = false,
  initialIsIntersecting = false,
  onChange,
}: UseIntersectionObserverOptions = {}): IntersectionReturn {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [ref, setRef] = useState<any>(null)

  const [state, setState] = useState<State>(() => ({
    isIntersecting: initialIsIntersecting,
    entry: undefined,
  }))

  const callbackRef = useRef<UseIntersectionObserverOptions['onChange']>(null)

  callbackRef.current = onChange

  const frozen = state.entry?.isIntersecting && freezeOnceVisible

  useEffect(() => {
    // Ensure we have a ref to observe
    if (!ref) return

    // Ensure the browser supports the Intersection Observer API
    if (!('IntersectionObserver' in window)) return

    // Skip if frozen
    if (frozen) return

    let unobserve: (() => void) | undefined

    const observer = new IntersectionObserver(
      (entries: IntersectionObserverEntry[]): void => {
        const thresholds = Array.isArray(observer.thresholds)
          ? observer.thresholds
          : [observer.thresholds]

        entries.forEach(entry => {
          const isIntersecting =
            entry.isIntersecting &&
            thresholds.some(threshold => entry.intersectionRatio >= threshold)

          setState({ isIntersecting, entry })

          if (callbackRef.current) {
            callbackRef.current(isIntersecting, entry)
          }

          if (isIntersecting && freezeOnceVisible && unobserve) {
            unobserve()
            unobserve = undefined
          }
        })
      },
      { threshold, root, rootMargin },
    )

    observer.observe(ref)

    return () => {
      observer.disconnect()
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    ref,
    // eslint-disable-next-line react-hooks/exhaustive-deps
    JSON.stringify(threshold),
    root,
    rootMargin,
    frozen,
    freezeOnceVisible,
  ])

  // ensures that if the observed element changes, the intersection observer is reinitialized
  const prevRef = useRef<Element | null>(null)

  useEffect(() => {
    if (
      !ref &&
      state.entry?.target &&
      !freezeOnceVisible &&
      !frozen &&
      prevRef.current !== state.entry.target
    ) {
      prevRef.current = state.entry.target
      setState({ isIntersecting: initialIsIntersecting, entry: undefined })
    }
  }, [ref, state.entry, freezeOnceVisible, frozen, initialIsIntersecting])

  const result = [
    setRef,
    ref,
    !!state.isIntersecting,
    state.entry,
  ] as IntersectionReturn

  // Support object destructuring, by adding the specific values.
  result.setRef = result[0]
  result.ref = result[1]
  result.isIntersecting = result[2]
  result.entry = result[3]

  return result
}