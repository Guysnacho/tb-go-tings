# Tigerbeetle TestBench

I wanna learn go and tigerbeetle so why not do it in one shot.

## Gameplan

I want a working bank basically. Lmfao, CUT AND DRY, SIMPLE SHIT. What's the worst that can happen?

A good approach would prolly be to implement all the functions on the [tigerbeetle reference](https://docs.tigerbeetle.com/reference).
Which basically handles the core functionality of a bank.
- create account
- transfer money
- query accounts
- check transaction history
- check balances

We'll add access control and allat too. Since its not postgres we don't get any RLS but that's okaaayyy. Also want to run a simulation like what Tigerbeetle already has on their site but for myself.

### The Heart

- [ ] Go api that serves all important operations. 
  - HTTP Handlers
    - [x] Hello world
    - [ ] CRUD for user (i.e. save transaction, fetch transactions, void transactions)
    - [ ] CRUD for admin (i.e. fetch transactions, void transactions)

### The Veins

We're currently publishing and consuming to 1 service.... since this is a testbench, I say that's fine. If we wanna break it up we'll do that later.

- [ ] Kafka messaging for non-save/modify operations altho it looks like tigerbeetle keeps most things immutable
  - topics for
    - [x] Hello world
    - [ ] creating transactions
    - [ ] creating accounts

### The Face

- [ ] Some sorta go library that isn't react for the frontend once we get there
  - bridges will be crossed once I get there.

## Project Structure

Figuring that out as we go but hoping it'll be self explanatory.
