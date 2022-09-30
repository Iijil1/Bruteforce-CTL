# Bruteforce CTL

One way to proof that a given turing machine can't halt is to give a Closed Tape Language for that machine. A detailed explanation can be found at https://www.sligocki.com/2022/06/10/ctl.html, but in short we want to find a set of configurations, described via regular expressions, such that the starting configuration is in there and every machine step from a configuration stays in the set.

The hard part about applying CTL is finding such a closed set of configurations. One can start by simulating the machine and adding the configurations one encounters to the set. But at some point abstractions need to be used or we end up with an infinite number of descriptions. The question then becomes at what point should a tape of '010101' be turned into '(01)+' or if that would help at all for this specific machine.

I attempt to get around having to answer that question by trying every abstraction possible, until I find one that works. That is done using two major components: the abstract tape representation and the decider itself.

# Abstract Tape Stack

The tape is split into two stacks, one each for the parts left and right of the head. Those tape stacks can contain abstractions, like '(01)+' or '(11|00)\*'.

We can push a symbol onto the stack which gives as a list of options for the resulting stack. One entry in the list would be the previos stack with the new symbol on top. The other entries are any abstractions that are found that can be created with the new symbol. For example pushing '1' onto '101101' might return '1011011' '10110(11)+' and '1(011)+'.

When we pop a symbol from the stack we get a list of symbols and the corresponding remaining stacks. The entries in the list are the results of the different possibilities to specify an abstraction at the top of the stack. For example popping '10(11|00)\*' would result in '1' + '0', '10(11|00)\*1' + '1' and '10(11|00)\*0' + '0'. 

The stack can handle a lot of abstractions, but using all of them can result in a searchspace that is too big for the decider to handle. So we can select options for which abstractions are allowed. For example we could specify to only try (A)+ after seeing A appearing 3 times in a row, or disable the use of (A|B)\* entirely.

# The Decider

The decider explores the space of configurations for a turing machine. Beginning with just the starting configuration we add more configurations by popping a symbol form the appropriate tape stack, applying the turing machine step on each of the results and pushing the resulting symbol onto the correct tape stack. We pick one of the results from the push and go with it.

Most likely we will eventually run into a situation where the turing machine halts. In that case we mark the configuration that lead to the halting state as unsuccessful and redo the predecessors. If we are able to select a different option for the push result that led to the halting configuration, then we continue with that. But if all possible options lead to halting then we mark this configuration as halting as well and redo its predecessors as well. That way we remove the entire branch of configurations that contain the abstraction that ultimately leads to halting.

If the starting configuration is ever marked as halting, then there are no possible push options left to explore and all the abstractions that were considered lead to halting. In that case we can't decide anything about the machine. Alternatively the process results in a closed set of configurations that includes the starting state. In that case we have proof that the machine doesn't halt.

# Iterative Deepening

Without an artificial limit to the search space the decider tends to try and make a too abstract left tape work by infinitely adding onto the right tape or vice versa. To prevent this we start off by limiting the depth of the search. Any configuration that takes too long to reach from the starting state is marked as halting due to the depth limit. When we find that the starting configuration itself halts due to that we increase the limit and start exploring again.

The configurations that were previously marked as halting due to the limit will be checked again. Information about configurations that halt for other reasons or that are known to be fine is kept and helps to make this deeper run faster than starting from scratch.