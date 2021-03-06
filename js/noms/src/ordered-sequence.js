// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

// @flow

import type {AsyncIterator, AsyncIteratorResult} from './async-iterator.js';
import {OrderedKey} from './sequence.js';
import type Value from './value.js'; // eslint-disable-line no-unused-vars
import {invariant, notNull} from './assert.js';
import search from './binary-search.js';
import Sequence, {SequenceCursor} from './sequence.js';

// See newCursorAt().
export function newCursorAtValue<T>(sequence: Sequence<T>, val: ?Value,
    forInsertion: boolean = false, last: boolean = false, readAhead: boolean = false)
    : Promise<OrderedSequenceCursor<any, any>> {
  let key;
  if (val !== null && val !== undefined) {
    key = new OrderedKey(val);
  }
  return newCursorAt(sequence, key, forInsertion, last, readAhead);
}


// Returns:
//   -null, if sequence is empty.
//   -null, if all values in sequence are < key.
//   -cursor positioned at
//      -first value, if |key| is null
//      -first value >= |key|
export async function newCursorAt<T>(sequence: ?Sequence<T>,
    key: ?OrderedKey<any>, forInsertion: boolean = false, last: boolean = false,
    readAhead: boolean = false): Promise<OrderedSequenceCursor<any, any>> {
  let cursor: ?OrderedSequenceCursor<any, any> = null;

  while (sequence) {
    cursor = new OrderedSequenceCursor(cursor, sequence, last ? -1 : 0, readAhead);
    if (key !== null && key !== undefined) {
      const lastPositionIfNotfound = forInsertion && sequence.isMeta;
      if (!cursor._seekTo(key, lastPositionIfNotfound)) {
        return cursor; // invalid
      }
    }

    sequence = await cursor.getChildSequence();
  }

  return notNull(cursor);
}

export class OrderedSequenceCursor<T, K: Value> extends
    SequenceCursor<T, Sequence<any>> {
  getCurrentKey(): OrderedKey<any> {
    invariant(this.idx >= 0 && this.idx < this.length);
    return this.sequence.getKey(this.idx);
  }

  clone(): OrderedSequenceCursor<T, K> {
    return new OrderedSequenceCursor(this.parent && this.parent.clone(), this.sequence, this.idx,
                                     this.readAhead);
  }

  // Moves the cursor to the first value in sequence >= key and returns true.
  // If none exists, returns false.
  _seekTo(key: OrderedKey<any>, lastPositionIfNotfound: boolean = false): boolean {
    // Find smallest idx where key(idx) >= key
    this.idx = search(this.length, i => this.sequence.getKey(i).compare(key) >= 0);

    if (this.idx === this.length && lastPositionIfNotfound) {
      invariant(this.idx > 0);
      this.idx--;
    }

    return this.idx < this.length;
  }
}

export class OrderedSequenceIterator<T, K: Value> { // AsyncIterator
  _iterator: Promise<AsyncIterator<T>>;

  constructor(cursorP: Promise<OrderedSequenceCursor<T, K>>) {
    this._iterator = cursorP.then(cur => cur.iterator());
  }

  next(): Promise<AsyncIteratorResult<T>> {
    return this._iterator.then(it => it.next());
  }

  return(): Promise<AsyncIteratorResult<T>> {
    return this._iterator.then(it => it.return());
  }
}
