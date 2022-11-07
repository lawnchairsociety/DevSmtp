namespace DevSmtp.Core.Commands
{
    public sealed class TurnResult : CommandResult
    {
        public TurnResult()
        {
        }

        public TurnResult(Exception error)
            : base(error)
        {
        }
    }
}
