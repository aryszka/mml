// Mikkamakka Lingua (MML) Syntax

// Notes:
// - this syntax definition distinguishes between expressions and statements. Expressions are statements with a
// return value. For simplicity, some composite expressions do not always have a return value but are handled as
// expressions in the syntax. These expressions are marked as 'fake' expressions, and the decision, whether they
// are really expressions, is made in a later step after parsing

mml = [ comment ] [ documentSequence ] { nlc } .

// - every token can be surrounded by comments
// - comments preceding a type alias, a type constraint or a definition by less than 2 new lines, and nothing
// else between than non-new-line whitespace, belong to the type alias, type constraint or definition
comment = "//" { CHARACTER } | "/*" { CHARACTER | nl } "*/" .

// - newline is handled as whitespace, except when it is between statements
// - statements in a sequence must by separated by ";" or newline
nl  = "\n" .
nlc = nl | comment .
sep = { nlc } ( nl | semicolon ) .

documentModeSwitch = modeSwitchWords
                     openBrace
                     documentSequence
                     closeBrace .

documentDelimitedStatement  = test | documentModeSwitch .
documentDelimitedStatements = documentDelimitedStatement { documentDelimitedStatement } .

documentStatements = ( statement | exportDefinition ) { sep ( statement | exportDefinition ) } [ sep ] .
documentSequence   = documentStatements | documentDelimitedStatements
                   | { documentStatements | documentDelimitedStatements } .

aliasWord           = { nlc } "alias" .
andWord             = { nlc } "and" .
breakWord           = { nlc } "break" .
caseWord            = { nlc } "case" .
complexRationalWord = { nlc } "complexRational" .
complexWord         = { nlc } "complex" .
continueWord        = { nlc } "continue" .
defaultWord         = { nlc } "default" .
deferWord           = { nlc } "defer" .
elseWord            = { nlc } "else" .
exportWord          = { nlc } "export" .
floatWord           = { nlc } "float" .
fnWord              = { nlc } "fn" .
forWord             = { nlc } "for" .
goWord              = { nlc } "go" .
ifWord              = { nlc } "if" .
importWord          = { nlc } "import" .
intWord             = { nlc } "int" .
inWord              = { nlc } "in" .
letWord             = { nlc } "let" .
orWord              = { nlc } "or" .
panicWord           = { nlc } "panic" .
rationalWord        = { nlc } "rational" .
receiveWord         = { nlc } "receive" .
recoverWord         = { nlc } "recover" .
sendWord            = { nlc } "send" .
switchWord          = { nlc } "switch" .
symbolWord          = { nlc } "symbol" .
testWord            = { nlc } "test" .
typeWord            = { nlc } "type" .

openParen   = { nlc } "(" .
closeParen  = { nlc } ")" .
openSquare  = { nlc } "[" .
closeSquare = { nlc } "]" .
openBrace   = { nlc } "{" .
closeBrace  = { nlc } "}" .

colon       = { nlc } ":" .
comma       = { nlc } "," .
communicate = { nlc } "<-" .
dot         = { nlc } "." .
question    = { nlc } "?" .
singleEq    = { nlc } "=" .
spread      = { nlc } "..." .
tilde       = { nlc } "~" .
semicolon   = { nlc } ";" .

andNot      = { nlc } "&~" .
diff        = { nlc } "-" .
div         = { nlc } "/" .
doubleAnd   = { nlc } "&&" .
doubleEq    = { nlc } "==" .
doubleOr    = { nlc } "||" .
greater     = { nlc } ">" .
greaterOrEq = { nlc } ">=" .
less        = { nlc } "<" .
lessOrEq    = { nlc } "<=" .
mod         = { nlc } "%" .
mul         = { nlc } "*" .
not         = { nlc } "!" .
notEq       = { nlc } "!=" .
power       = { nlc } "^" .
shiftLeft   = { nlc } "<<" .
shiftRight  = { nlc } ">>" .
singleAnd   = { nlc } "&" .
singleOr    = { nlc } "|" .
sum         = { nlc } "+" .
incOne      = { nlc } "++" .
decOne      = { nlc } "--" .

int             = { nlc } DIGIT { DIGIT } . // or something like that, initially
float           = { nlc } DIGIT dot { DIGIT } .
rational        = { nlc } rationalWord openParen int comma int closeParen .
complex         = { nlc } complexWord openParen float comma float closeParen .
complexRational = { nlc } complexRationalWord rational comma rational closeParen .
string          = { nlc } """" { CHARACTER } """" . // in the future, strings can contain new-lines
bool            = { nlc } ( "true" | "false" ) .
channel         = { nlc } "<>" .
symbol          = { nlc } LETTER { LETTER | DIGIT } .

staticSymbol     = symbol | string .
dynamicSymbol    = symbolWord openParen expression closeParen .
symbolExpression = staticSymbol | dynamicSymbol .

listExpressionSequence = ( expression | spreadExpression )
                         { comma ( expression | spreadExpression ) }
                         [ comma ] .
list                   = openSquare [ listExpressionSequence ] closeSquare .
mutableList            = tilde openSquare [ listExpressionSequence ] closeSquare .

mutableStructure = tilde structure .
structureItem    = symbolExpression colon expression .
structure        = openBrace [ structureItem { comma structureItem } [ comma ] ] closeBrace .

andExpression = andWord openParen [ expressionSequence ] closeParen .
orExpression  = orWord openParen [ expressionSequence ] closeParen .

arithmeticOperator = andNot
                   | diff
                   | div
                   | mod
                   | mul
                   | power
                   | shiftLeft
                   | shiftRight
                   | singleAnd
                   | singleOr
                   | sum
                   | tilde .

binaryOperator = arithmeticOperator
               | channelOp
               | doubleAnd
               | doubleEq
               | doubleOr
               | greater
               | greaterOrEq
               | less
               | lessOrEq
               | notEq .

unaryOperator = diff
              | not
              | sum
              | tilde .

binaryOperation = expression binaryOperator expression .
unaryOperation = unaryOperator expression .

test = testWord openBrace statementSequence closeBrace .

// is this really needed or is it possible to figure the return type of the basic operators?
// or is it possible to enforce that the operands be the same? e.g:
// type sum fn(a int|float|rational|string) a
modeSwitchWords = ( intWord | floatWord | rationalWord | complexWord | complexRationalWord ) .
modeSwitch      = modeSwitchWords
                  openBrace
                  statementSequence
                  closeBrace .

limitedStatement = test | modeSwitch .
limitedStatements = limitedStatement { limitedStatement } .

statements        = statement { sep statement } [ sep ] .
statementSequence = { statements | limitedStatements } .

functionBody      = openBrace statementSequence closeBrace .
spreadSymbol      = spread staticSymbol .

functionFact = openParen
               ( [ staticSymbol { comma staticSymbol } [ comma spreadSymbol ] [ comma ] ]
               | spreadSymbol [ comma ] )
               closeParen
               ( statement | functionBody | test ) .

function         = fnWord functionFact .
mutatingFunction = fnWord tilde functionFact .

rangeExpression = [ expression ] colon [ expression ] .
expressionQuery = expression openSquare ( expression | rangeExpression ) closeSquare .
symbolQuery     = expression dot symbolExpression .
query           = expressionQuery | symbolQuery .

spreadExpression = spread expression .

// fake expression:
functionCall = [ goWord | deferWord ] expression openParen
               ( [ expression { comma expression } [ comma spreadExpression ] [ comma ] ]
               | spreadExpression [ comma ] )
               closeParen .

panicCall   = panicWord openParen expression closeParen .
recoverCall = [ deferWord ] recoverWord openParen closeParen .

receive = receiveWord openParen expression closeParen
        | channelOp expression .
send    = sendWord openParen expression comma expression closeParen
        | expression channelOp expression .

assignment = ( staticSymbol | query ) [ arithmeticOperator ] singleEq expression
           | staticSymbol incOne
           | staticSymbol decOne .

typeMatch = typeReference doubleEq typeExpression
            { comma typeReference doubleEq typeExpression }
            [ comma ]
          | typeWord openParen
            staticSymbol doubleEq typeExpression
            { comma staticSymbol doubleEq typeExpression }
            [ comma ]
            closeParen

queryMatch         = staticSymbol [ singleEq ] query { comma staticSymbol [ singleEq ] query } .
communicationMatch = [ staticSymbol singleEq ] receive | send .
matchExpression    = expression | queryMatch | typeMatch | communicationMatch .

switchClause      = caseWord matchExpression colon statementSequence .
switchValueClause = caseWord expression colon statementSequence .
swtichTypeClause  = caseWord typeExpression colon statementSequence .
defaultClause     = defaultWord colon statementSequence .

switchConditional = switchWord openBrace
                    { switchClause }
                    defaultClause
                    { switchClause }
                    closeBrace .

switchValueConditional = switchWord expression openBrace
                         { switchValueClause }
                         defaultClause
                         { switchValueClause }
                         closeBrace .

switchTypeConditional = switchWord typeWord expression openBrace
                        { switchTypeClause } 
                        [ defaultClause ]
                        { switchTypeClause }
                        closeBrace .

questionConditional = matchExpression question statement colon statement .

ifClause      = openBrace statementSequence closeBrace .
ifCondition   = ifWord matchExpression ifClause .
ifConditional = ifCondition [ elseWord ifCondition ] elseWord ifClause .

// fake expression:
conditional = switchConditional
            | switchValueConditional
            | switchTypeConditional
            | ifConditional
            | questionConditional .

voidIfConditional          = ifCondition .
voidSwitchConditional      = switchWord openBrace { swtichClause } closeBrace .
voidSwitchValueConditional = switchWord expression openBrace { switchValueClause } closeBrace .
voidSwitchTypeConditional  = switchWord typeWord expression openBrace { switchTypeClause } closeBrace .

voidConditional = voidSwitchConditional
                | voidSwitchValueConditional
                | voidSwitchTypeConditional
                | voidIfConditional .

expression = int
           | float
           | rational
           | complex
           | complexRational
           | string
           | bool
           | symbol
	   | dynamicSymbol
           | list
           | mutableList
           | structure
           | mutableStructure
           | channel
           | andExpression
           | orExpression
           | unaryOperation
           | binaryOperation
           | function
           | mutatingFunction
           | query
           | functionCall
           | panicCall
           | recoverCall
           | receive
           | conditional
           | importExpression
           | expressionGroup .

expressionGroup    = openParen expression closeParen .
expressionSequence = expression { comma expression } [ comma ] .

loopControl     = breakWord | continueWord .
conditionalLoop = forWord matchExpression statement .
memberLoop      = forWord staticSymbol inWord expression statement .
rangeLoop       = forWord staticSymbol inWord rangeExpression statement .
loop            = conditionalLoop | memberLoop | rangeLoop .

definitionItem        = staticSymbol [ singleEq ] expression .
mutableDefinitionItem = tilde definitionItem .
valueDefinition       = letWord ( definitionItem | mutableDefinitionItem ) .

valueDefinitionGroup = letWord
                       openParen
                       [ definitionItem | mutableDefinitionItem
                         { comma ( definitionItem | mutableDefinitionItem ) }
                         [ comma ] ]
                       closeParen .

functionDefinitionItem = staticSymbol functionFact .

mutatingFunctionDefinitionItem = tilde functionDefinitionItem .
functionDefinition             = fnWord ( functionDefinitionItem | mutatingFunctionDefinitionItem ) .

functionDefinitionGroup = fnWord openParen
                          [ functionDefinitionItem | mutatingFunctionDefinitionItem
                            { comma functionDefinitionItem | comma mutatingFunctionDefinitionItem }
                            [ comma ] ]
                          closeParen .

definition = valueDefinition
           | valueDefinitionGroup
           | functionDefinition
           | functionDefinitionGroup .

primitiveType = intWord
              | floatWord
              | rationalWord
              | complexWord
              | complexRationalWord
              | stringWord
              | boolWord .

listType        = openSquare typeExpression closeSquare .
mutableListType = tilde openSquare typeExpression closeSquare .

structureTypeItem    = staticSymbol colon typeExpression .
structureType        = openBrace [ structureTypeItem { comma structureTypeItem } [ comma ] ] closeBrace .
mutableStructureType = tilde structureType .
// this basically means that that every value may pass the type constraint of {}
// this may not stay like that forever, it is not the goal to introduce a virtual nil type
// (and definitely not a nil value)

argumentType = [ staticSymbol ] typeExpression .
functionType = fnWord openParen
               ( [ argumentType { comma argumentType } [ comma spread argumentType ] [ comma ] ] |
                 spread argumentType [ comma ] )
               closeParen
               [ argumentType ] .

channelType   = less typeExpression greater .
typeReference = typeWord staticSymbol .

unionType = typeExpression { singleOr typeExpression } .

typeExpression = primitiveType
               | listType
               | mutableListType
               | structureType
               | mutableStructureType
               | functionType
               | channelType
               | unionType
               | typeReference .

typeConstraint = typeReference [ singleEq ] typeExpression .
typeAlias      = typeWord aliasWord staticSymbol [ singleEq ] typeExpression .

importExpression      = importWord string .
importItem            = staticSymbol [ singleEq ] string .
importDefinition      = importWord importItem .
importDefinitionGroup = importWord openParen [ importItem { comma importItem } [ comma ] ] closeParen .

definition = valueDefinition
           | valueDefinitionGroup
           | functionDefinition
           | functionDefinitionGroup
           | importDefinition
           | importDefinitionGroup .

exportDefinition = exportWord
                   ( definition
                   | staticSymbol
                   | openParen [ staticSymbol { comma staticSymbol } [ comma ] ] closeParen
                   | importExpression ) .

statement = expression
          | assignment
          | send
          | close
          | voidConditional
          | loopControl
          | loop
          | definition
          | typeConstraint
          | typeAlias
          | importStatement .
