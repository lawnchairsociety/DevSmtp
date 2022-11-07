namespace DevSmtp.Core.Commands
{
    public sealed class DataResult : CommandResult
    {
        public DataResult()
        {
        }

        public DataResult(Exception error)
            : base(error)
        {
        }
    }
}
