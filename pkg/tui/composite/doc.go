package composite

// This package is responsible for creating an intermediary layer between runtime data and information to be rendered.
// This ensures that there's a decoupling between runtime state and render state representations
// Also, it provides a natural space for encoding the information before representing it, so things like
// shading/decoriation can be incorporated. Finally, since the block size is bigger than the board size,
// The cells can create transition effects between each section.
//
// A quick definition of the domain here is that a cell is a single unit of visual information and a group
// of cells together is a block.
