namespace DevSmtp.Core.Commands
{
    public sealed class QuitResult : CommandResult
    {
        public QuitResult()
        {
        }

        public QuitResult(Exception error)
            : base(error)
        {
        }
    }
}
