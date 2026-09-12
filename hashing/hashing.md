---
title: "Day 1 — Hash Functions and the Foundation of Integrity"
date: 2026-09-12
series: "Building a Blockchain from Scratch in Go"
week: 1
day: 1
tags: [blockchain, go, cryptography, sha256]
---

# Day 1 — Hash Functions and the Foundation of Integrity

## What a hash function is

A hash function maps input of any length to output of a fixed length. SHA-256 takes
anything — a single byte, a book, a 4 GB video — and returns exactly 32 bytes.

It has three security properties, and they are not equally strong:

| Property | Meaning | Attack cost |
|---|---|---|
| Preimage resistance | Given `h`, you cannot find `m` such that `H(m) = h` | 2²⁵⁶ |
| Second-preimage resistance | Given `m₁`, you cannot find `m₂ ≠ m₁` with `H(m₂) = H(m₁)` | 2²⁵⁶ |
| Collision resistance | You cannot find *any* pair `m₁ ≠ m₂` with `H(m₁) = H(m₂)` | **2¹²⁸** |

The third row is the one worth remembering. Finding *some* collision is vastly
cheaper than forging a *specific* target.

The reason is the birthday paradox. Among random inputs, a collision appears after
roughly 2^(n/2) attempts rather than 2^n. So collision resistance is always the
weakest of the three, and a protocol has to assume the attacker's cheapest path
rather than its most expensive one.

## Hashing is not encryption

This distinction gets muddled constantly:

- **Encryption is reversible.** With the key, you recover the original.
- **Hashing is irreversible.** There is no key, and no amount of cleverness inverts it.

The reason is information-theoretic, not computational. SHA-256 compresses arbitrary
input into 32 bytes, so information is necessarily destroyed. Asking to "decrypt a
hash" is asking to recover a 1 GB file from 32 bytes. The question is malformed.

## The avalanche effect

Hash four nearly identical strings — one differing by a case change, one by an added
period, one by an extra space — and the four outputs share no visible structure.
No common prefix, no similar pattern. They look like four unrelated random numbers.

This is the avalanche effect: flipping one input bit changes roughly half of all
output bits.

Why it matters: there is no "close enough." You cannot nudge a record slightly and
get a slightly different hash. Any modification, however small, produces a completely
unrelated value. Tampering isn't merely detectable — it's loudly detectable.

### A case that surprised me

Two strings that render identically in a terminal can have different hashes, if one
contains a homoglyph — say U+0435, the Cyrillic small letter *ye*, in place of the
Latin `e`.

Different bytes, different hash. The function is doing exactly the right thing: it
operates on bytes, not on glyphs, and has no notion of what "looks the same" means.

Filing this away as a real concern for later. If content can contain homoglyphs, two
records that appear identical to a human will have different hashes and different
identifiers. Not a cryptographic problem, but a deduplication and UI problem.

## What a hash can and cannot prove

This is the most important idea of the day.

A hash **can** prove that data has not changed since the hash was computed. Recompute,
compare, and any difference shows up immediately.

A hash **cannot** prove that the data was ever correct.

Write down "the Moon is made of cheese" and hash it. The hash is perfectly valid. It
proves nobody edited the false statement. It says nothing about the Moon.

> **Integrity is not truth.**

This is a hard architectural constraint, not a philosophical aside. A tamper-evident
record of a wrong fact is still a wrong fact. Any system built on hash chains has to
be honest about which of those two things it actually guarantees.

## Reading: the Bitcoin whitepaper

Read the Abstract, §1 (Introduction), and §2 (Transactions). About 20 minutes.
Official PDF: https://bitcoin.org/bitcoin.pdf

Deliberately skipping the rest for now — proof-of-work and the network section are
week 6–8 material and would mostly be noise at this point.

### From §1

The paper frames the problem as one of trust rather than technology. Relying on a
trusted intermediary imposes costs beyond fees: limits on transaction size, the
practical impossibility of non-reversible payments, and a need to collect customer
information that wouldn't otherwise be required.

What struck me is what the paper does **not** claim. There's no argument anywhere that
an immutable record is therefore a true record. The whole document is about
establishing *order* and *authorization*, never *veracity*. The original design was
always narrower than the popular understanding of it.

### From §2 — the pivotal idea

Two things here reframed my mental model.

**A coin is not a balance.** It's defined as a chain of digital signatures, where each
transfer signs over the previous transaction plus the next owner's public key. There's
no account holding a number. Ownership is a verifiable history of transfers.

**Signatures alone cannot solve double-spending.** This took a while to sit with.

Suppose I own one coin and sign two transfers of it — one to Alice, one to Bob. Both
signatures are cryptographically valid. Both prove I authorized that transfer. Examined
individually, neither is detectably fraudulent.

Signatures answer *"did the owner authorize this?"* They cannot answer *"was this the
first transfer of this coin?"* — because that question isn't about any single
transaction. It's about the relationship between transactions, which makes it a
question about **order**.

So preventing double-spending requires everyone to agree on which transfer came first.
Not most people. Everyone, with no trusted party to ask.

**That is the entire reason consensus exists.** Not to prevent tampering — the hash
chain does that. Consensus exists to establish a single agreed-upon ordering among
mutually distrusting parties.

## Three problems, three mechanisms

This is what actually clicked today. I'd been treating "blockchain" as one thing.
It's three independent mechanisms solving three separate problems:

| Problem | Mechanism |
|---|---|
| Was this record modified after the fact? | Hash chain |
| Did the claimed author really authorize this? | Digital signature |
| Which of two conflicting records came first? | Consensus |

A hash chain with no replication is just a log I can rewrite at will. Signatures
without ordering can't stop double-spending. Consensus without hashing has nothing
to anchor. Each one is necessary and none is sufficient.

## Questions I can now answer

**What's the difference between hashing and encryption?**
Encryption is reversible with a key; hashing is irreversible and keyless. They solve
different problems — confidentiality versus integrity.

**Why can't SHA-256 be decrypted?**
It compresses arbitrary-length input to 32 bytes, destroying information by
construction. The inverse function doesn't exist. This is information theory, not
insufficient computing power.

**Can a hash prove the input is authentic?**
No. It proves the input hasn't changed since hashing. Authenticity of the *author*
needs a signature. Correctness of the *content* needs something outside cryptography
entirely.

**Why is collision resistance the weakest property?**
The birthday paradox — finding any colliding pair takes ~2^(n/2) attempts, not 2^n.

## Questions I still can't answer

- How do multiple machines agree on the next block?
- How does one block hold many transactions?
- How can a client verify a record is on-chain without downloading everything?
- Why does the hash *input* need to be a uniquely determined byte sequence? Apparently
  this is the single most dangerous bug class in the whole system.

## Tomorrow

Quantify what I only observed today. Two experiments:

1. **Measure the avalanche effect** — flip every input bit in turn, count how many
   output bits change, confirm the average sits near 50%.
2. **Produce an actual hash collision** — truncate SHA-256 to 32 bits and use a
   birthday attack. Should take seconds. Turning `2^(n/2)` from a formula into
   something I've personally watched happen.